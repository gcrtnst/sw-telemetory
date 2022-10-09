function test()
    local test_tbl = {
        {"testClientSizeNormal", testClientSizeNormal},
        {"testClientSizeError", testClientSizeError},
        {"testClientBusyAfterInit", testClientBusyAfterInit},
        {"testClientBusyAfterGet", testClientBusyAfterGet},
        {"testClientBusyAfterCancel", testClientBusyAfterCancel},
        {"testClientBusyAfterTimeout", testClientBusyAfterTimeout},
        {"testClientBusyAfterReply", testClientBusyAfterReply},
        {"testClientBusyAfterCancelTimeout", testClientBusyAfterCancelTimeout},
        {"testClientBusyAfterCancelReply", testClientBusyAfterCancelReply},
        {"testEncodeCSVRecord", testEncodeCSVRecord},
        {"testEncodeCSVField", testEncodeCSVField},
        {"testEscapeQuery", testEscapeQuery},
    }

    local t = buildT()
    for _, test_entry in ipairs(test_tbl) do
        local test_name, test_fn = table.unpack(test_entry)
        t:reset()

        local is_success, err = pcall(test_fn, t)
        if not is_success then
            io.write(string.format("FAIL %s\n", test_name))
            io.write(string.format("     %s\n", err))
        else
            io.write(string.format("PASS %s\n", test_name))
        end
    end
end

function testClientSizeNormal(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/ur", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(1, 52149, "/ur")
end

function testClientSizeError(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_call("ctx", t.env.c_client_status_size, nil)
    t.env.async._assert_cnt(0)
end

function testClientBusyAfterInit(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(1, 52149, "/url")
end

function testClientBusyAfterGet(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_call("ctx", t.env.c_client_status_busy, nil)
    t.env.async._assert_cnt(1)
end

function testClientBusyAfterCancel(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_call("ctx", t.env.c_client_status_busy, nil)
    t.env.async._assert_cnt(1)
end

function testClientBusyAfterTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientOnTick()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(2, 52149, "/url")
end

function testClientBusyAfterReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpReply(52149, "/url", "resp")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(2, 52149, "/url")
end

function testClientBusyAfterCancelTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()
    t.env.clientOnTick()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(2, 52149, "/url")
end

function testClientBusyAfterCancelReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback()

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()
    t.env.clientHttpReply(52149, "/url", "resp")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback.assert_wait()
    t.env.async._assert_call(2, 52149, "/url")
end

function testEncodeCSVRecord(t)
    local tests = {
        {nil, nil},
        {{0}, nil},
        {{}, "\r\n"},
        {{""}, "\r\n"},
        {{"", ""}, ",\r\n"},
        {{"", "", ""}, ",,\r\n"},
        {{"a"}, "a\r\n"},
        {{"a", "a"}, "a,a\r\n"},
        {{"a", "a", "a"}, "a,a,a\r\n"},
        {{"", "", "a"}, ",,a\r\n"},
        {{"", "a", ""}, ",a,\r\n"},
        {{"", "a", "a"}, ",a,a\r\n"},
        {{"a", "", ""}, "a,,\r\n"},
        {{"a", "", "a"}, "a,,a\r\n"},
        {{"a", "a", ""}, "a,a,\r\n"},
        {{"a", "a", "a"}, "a,a,a\r\n"},
        {{"abc", "def", "ghi"}, "abc,def,ghi\r\n"},
        {{",", ",", ","}, '",",",",","\r\n'},
    }

    for i, tt in ipairs(tests) do
        local in_record, want_s = table.unpack(tt)
        t:reset()
        t.fn()

        local got_s = t.env.encodeCSVRecord(in_record)
        if got_s ~= want_s then
            error(string.format('case %d: expected "%s", got "%s"', i, want_s, got_s))
        end
    end
end

function testEncodeCSVField(t)
    local tests = {
        {nil, nil},
        {'', ''},
        {'a', 'a'},
        {' ', ' '},
        {'\r', '\r'},
        {'\n', '\n'},
        {'\r\n', '"\r\n"'},
        {'"', '""""'},
        {',', '","'},
        {'abc', 'abc'},
        {'"abc"', '"""abc"""'},
        {'a"b', '"a""b"'},
        {'"a"b"', '"""a""b"""'},
        {' abc', ' abc'},
        {'abc,def', '"abc,def"'},
        {'abc\ndef', 'abc\ndef'},
        {'abc\rdef', 'abc\rdef'},
        {'abc\r\ndef', '"abc\r\ndef"'},
    }

    for i, tt in ipairs(tests) do
        local in_s, want_s = table.unpack(tt)
        t:reset()
        t.fn()

        local got_s = t.env.encodeCSVField(in_s)
        if got_s ~= want_s then
            error(string.format('case %d: expected "%s", got "%s"', i, want_s, got_s))
        end
    end
end

function testEscapeQuery(t)
    local tests = {
        {nil, nil},
        {"", ""},
        {"abc", "abc"},
        {"one two", "one+two"},
        {"10%", "10%25"},
        {" ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;", "+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B"},
    }

    for i, tt in ipairs(tests) do
        local in_s, want_s = table.unpack(tt)
        t:reset()
        t.fn()

        local got_s = t.env.escapeQuery(in_s)
        if got_s ~= want_s then
            error(string.format('case %d: expected "%s", got "%s"', i, want_s, got_s))
        end
    end
end

function buildT()
    local env = {}
    local fn, err = loadfile("script.lua", "t", env)
    if fn == nil then
        error(err)
    end

    local t = {
        env = env,
        fn = fn,
        reset = function(self)
            for k, _ in pairs(self.env) do
                self.env[k] = nil
            end
            self.env.ipairs = ipairs
            self.env.math = math
            self.env.next = next
            self.env.pairs = pairs
            self.env.string = string
            self.env.table = table
            self.env.tonumber = tonumber
            self.env.tostring = tostring
            self.env.type = type
            self.env.async = buildMockAsync()
        end
    }
    t:reset()
    return t
end

function buildMockAsync()
    local async = {
        _cnt = 0,
        _port = nil,
        _url = nil,
    }

    async.httpGet = function(port, url)
        async._cnt = async._cnt + 1
        async._port = port
        async._url = url
    end

    async._assert_cnt = function(cnt)
        if async._cnt ~= cnt then
            error(string.format("%q", async._cnt))
        end
    end

    async._assert_call = function(cnt, port, url)
        if async._cnt ~= cnt then
            error(string.format("%q", async._cnt))
        end
        if async._port ~= port then
            error(string.format("%q", async._port))
        end
        if async._url ~= url then
            error(string.format("%q", async._url))
        end
    end

    return async
end

function buildMockClientCallback()
    local callback = {
        cnt = 0,
        ctx = nil,
        status = nil,
        resp = nil,
    }

    callback.fn = function(ctx, status, resp)
        callback.cnt = callback.cnt + 1
        callback.ctx = ctx
        callback.status = status
        callback.resp = resp
    end

    callback.assert_wait = function()
        if callback.cnt ~= 0 then
            error(string.format("%q", callback.cnt))
        end
    end

    callback.assert_call = function(ctx, status, resp)
        if callback.cnt ~= 1 then
            error(string.format("%q", callback.cnt))
        end
        if callback.ctx ~= ctx then
            error(string.format("%q", callback.ctx))
        end
        if callback.status ~= status then
            error(string.format("%q", callback.status))
        end
        if callback.resp ~= resp then
            error(string.format("%q", callback.resp))
        end
    end

    return callback
end

test()
