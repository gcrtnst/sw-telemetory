g_test_tbl = {}

function test()
    local test_tbl = {}
    for test_name, test_fn in pairs(g_test_tbl) do
        table.insert(test_tbl, {test_name, test_fn})
    end
    table.sort(test_tbl, function(x, y)
        return x[1] < y[1]
    end)

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

function g_test_tbl.testClientSizeNormal(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/ur", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(1, 52149, "/ur")
end

function g_test_tbl.testClientSizeError(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_size, nil)
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testClientBusyAfterInit(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(1, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterGet(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_busy, nil)
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testClientBusyAfterCancel(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_busy, nil)
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testClientBusyAfterTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientOnTick()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpReply(52149, "/url", "resp")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterCancelTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()
    t.env.clientOnTick()

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterCancelReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()
    t.env.clientHttpReply(52149, "/url", "resp")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertWait()
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientCancel(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpCancel()
    callback:assertCall("ctx", t.env.c_client_status_cancel, nil)
end

function g_test_tbl.testClientCancelIdle(t)
    t:reset()
    t.fn()
    t.env.clientHttpCancel()
end

function g_test_tbl.testClientCancelCancel(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpCancel()
    t.env.clientHttpCancel()
    callback:assertCall("ctx", t.env.c_client_status_cancel, nil)
end

function g_test_tbl.testClientCancelTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpCancel()
    t.env.clientOnTick()
    callback:assertCall("ctx", t.env.c_client_status_cancel, nil)

    -- confirm timeout
    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientCancelReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpCancel()
    t.env.clientHttpReply(52149, "/url", "resp")
    callback:assertCall("ctx", t.env.c_client_status_cancel, nil)

    -- confirm reply
    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientTimeoutBefore(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 3
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientOnTick()
    t.env.clientOnTick()
    callback:assertWait()

    -- confirm busy
    local callback = buildMockClientCallback("callback")
    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_busy, nil)
end

function g_test_tbl.testClientTimeoutAfter(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 3
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientOnTick()
    t.env.clientOnTick()
    t.env.clientOnTick()
    callback:assertCall("ctx", t.env.c_client_status_timeout, nil)

    -- confirm idle
    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientTimeoutGet(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 3
    local callback_called = false
    local callback = function(ctx, status, resp)
        callback_called = true

        t.env.clientHttpGet("ctx", 52149, "/url", function() end)
        t.env.async:assertCall(2, 52149, "/url")
    end

    t.env.clientHttpGet("ctx", 52149, "/url", callback)
    t.env.clientOnTick()
    t.env.clientOnTick()
    t.env.clientOnTick()
    assertEqual("callback_called", true, callback_called)
end

function g_test_tbl.testClientReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpReply(52149, "/url", "resp")
    callback:assertCall("ctx", t.env.c_client_status_ok, "resp")

    -- confirm idle
    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientReplyIgnoreIdle(t)
    t:reset()
    t.fn()
    t.env.clientHttpReply(52149, "/url", "resp")
end

function g_test_tbl.testClientReplyIgnorePort(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpReply(52148, "/url", "resp")
    callback:assertWait()

    -- confirm busy
    local callback = buildMockClientCallback("callback")
    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_busy, nil)
end

function g_test_tbl.testClientReplyIgnoreReq(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    t.env.clientHttpReply(52149, "/dmy", "resp")
    callback:assertWait()

    -- confirm busy
    local callback = buildMockClientCallback("callback")
    t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    callback:assertCall("ctx", t.env.c_client_status_busy, nil)
end

function g_test_tbl.testClientReplyGet(t)
    t:reset()
    t.fn()

    local callback_called = false
    local callback = function(ctx, status, resp)
        callback_called = true

        t.env.clientHttpGet("ctx", 52149, "/url", function() end)
        t.env.async:assertCall(2, 52149, "/url")
    end

    t.env.clientHttpGet("ctx", 52149, "/url", callback)
    t.env.clientHttpReply(52149, "/url", "resp")
    assertEqual("callback_called", true, callback_called)
end

function g_test_tbl.testEncodeCSVRecord(t)
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
        assertEqual(string.format("case %d", i), want_s, got_s)
    end
end

function g_test_tbl.testEncodeCSVField(t)
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
        assertEqual(string.format("case %d", i), want_s, got_s)
    end
end

function g_test_tbl.testEscapeQuery(t)
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
        assertEqual(string.format("case %d", i), want_s, got_s)
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

    async.httpGet = function(...)
        return async:_httpGet(...)
    end

    async._httpGet = function(self, port, url)
        self._cnt = self._cnt + 1
        self._port = port
        self._url = url
    end

    async.assertCallCount = function(self, cnt)
        assertEqual("async._cnt", cnt, self._cnt)
    end

    async.assertCall = function(self, cnt, port, url)
        assertEqual("async._cnt", cnt, self._cnt)
        assertEqual("async._port", port, self._port)
        assertEqual("async._url", url, self._url)
    end

    return async
end

function buildMockClientCallback(name)
    local callback = {
        name = name,
        cnt = 0,
        ctx = nil,
        status = nil,
        resp = nil,
    }

    callback.fn = function(...)
        return callback:_fn(...)
    end

    callback._fn = function(self, ctx, status, resp)
        self.cnt = self.cnt + 1
        self.ctx = ctx
        self.status = status
        self.resp = resp
    end

    callback.assertWait = function(self)
        assertEqual(string.format("%s.cnt", self.name), 0, self.cnt)
    end

    callback.assertCall = function(self, ctx, status, resp)
        assertEqual(string.format("%s.cnt", self.name), 1, self.cnt)
        assertEqual(string.format("%s.ctx", self.name), ctx, self.ctx)
        assertEqual(string.format("%s.status", self.name), status, self.status)
        assertEqual(string.format("%s.resp", self.name), resp, self.resp)
    end

    return callback
end

function assertEqual(name, want, got)
    if got ~= want then
        error(string.format("%s: expected %q, got %q", name, want, got))
    end
end

test()
