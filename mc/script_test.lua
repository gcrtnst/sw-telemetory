function test()
    local test_tbl = {
        {"testEscapeQuery", testEscapeQuery},
        {"testEncodeCSVRecord", testEncodeCSVRecord},
        {"testEncodeCSVField", testEncodeCSVField},
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

function testEscapeQuery(t)
    local tests = {
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

function testEncodeCSVRecord(t)
    local tests = {
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
        end
    }
    t:reset()
    return t
end

test()
