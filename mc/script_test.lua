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
    local s = "PASS"
    for _, test_entry in ipairs(test_tbl) do
        local test_name, test_fn = table.unpack(test_entry)
        t:reset()

        local is_success, err = pcall(test_fn, t)
        if not is_success then
            io.write(string.format("FAIL %s\n", test_name))
            io.write(string.format("     %s\n", err))
            s = "FAIL"
        else
            io.write(string.format("PASS %s\n", test_name))
        end
    end
    io.write(s .. "\n")
end

function g_test_tbl.testLogActiveNil(t)
    t:reset()
    t.env.c_log_active_ch = nil
    t.env.c_log_bool_ch_start = 1
    t.env.c_log_bool_ch_limit = 0
    t.env.c_log_number_ch_start = 1
    t.env.c_log_number_ch_limit = 0
    t.env.input._bool[32] = false
    t.env.property._number["Port"] = 52149
    t.env.property._text["Title"] = "title"
    t.fn()

    t.env.logOnTick()
    t.env.async:assertCall(1, 52149, "/time")
end

function g_test_tbl.testLogActiveFalse(t)
    t:reset()
    t.env.c_log_active_ch = 32
    t.env.c_log_bool_ch_start = 1
    t.env.c_log_bool_ch_limit = 0
    t.env.c_log_number_ch_start = 1
    t.env.c_log_number_ch_limit = 0
    t.env.input._bool[32] = false
    t.env.property._number["Port"] = 52149
    t.env.property._text["Title"] = "title"
    t.fn()

    t.env.logOnTick()
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testLogActiveTrue(t)
    t:reset()
    t.env.c_log_active_ch = 32
    t.env.c_log_bool_ch_start = 1
    t.env.c_log_bool_ch_limit = 0
    t.env.c_log_number_ch_start = 1
    t.env.c_log_number_ch_limit = 0
    t.env.input._bool[32] = true
    t.env.property._number["Port"] = 52149
    t.env.property._text["Title"] = "title"
    t.fn()

    t.env.logOnTick()
    t.env.async:assertCall(1, 52149, "/time")
end

function g_test_tbl.testLogCancel(t)
    t:reset()
    t.env.c_log_active_ch = 32
    t.env.c_log_bool_ch_start = 1
    t.env.c_log_bool_ch_limit = 0
    t.env.c_log_number_ch_start = 1
    t.env.c_log_number_ch_limit = 0
    t.env.input._bool[32] = true
    t.env.property._number["Port"] = 52149
    t.env.property._text["Title"] = "title"
    t.fn()

    t.env.logOnTick()
    t.env.async:assertCall(1, 52149, "/time")

    t.env.input._bool[32] = false
    t.env.logOnTick()

    -- confirm send cancel
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testLogHeader(t)
    local tests = {
        {
            1, 0, 1, 0,
            {},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%0D%0A0%0D%0A"
        },
        {
            1, 1, 1, 0,
            {["Bool Label 1"] = "BL1"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CBL1%0D%0A0%2CFALSE%0D%0A"
        },
        {
            2, 2, 1, 0,
            {["Bool Label 2"] = "BL2"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CBL2%0D%0A0%2CFALSE%0D%0A"
        },
        {
            1, 3, 1, 0,
            {["Bool Label 1"] = "BL1", ["Bool Label 2"] = "BL2", ["Bool Label 3"] = "BL3"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CBL1%2CBL2%2CBL3%0D%0A0%2CFALSE%2CFALSE%2CFALSE%0D%0A"
        },
        {
            1, 0, 1, 1,
            {["Number Label 1"] = "NL1"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CNL1%0D%0A0%2C0%0D%0A"
        },
        {
            1, 0, 2, 2,
            {["Number Label 2"] = "NL2"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CNL2%0D%0A0%2C0%0D%0A"
        },
        {
            1, 0, 1, 3,
            {["Number Label 1"] = "NL1", ["Number Label 2"] = "NL2", ["Number Label 3"] = "NL3"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CNL1%2CNL2%2CNL3%0D%0A0%2C0%2C0%2C0%0D%0A"
        },
        {
            1, 1, 1, 1,
            {["Bool Label 1"] = "BL1", ["Number Label 1"] = "NL1"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CBL1%2CNL1%0D%0A0%2CFALSE%2C0%0D%0A"
        },
        {
            3, 4, 1, 2,
            {["Bool Label 3"] = "BL3", ["Bool Label 4"] = "BL4", ["Number Label 1"] = "NL1", ["Number Label 2"] = "NL2"},
            "/write?path=title%2Ftitle-20060102150405.csv&data=%23%2CBL3%2CBL4%2CNL1%2CNL2%0D%0A0%2CFALSE%2CFALSE%2C0%2C0%0D%0A"
        },
    }

    for i, tt in ipairs(tests) do
        local in_c_log_bool_ch_start, in_c_log_bool_ch_limit, in_c_log_number_start, in_c_log_number_limit, in_property_text, want_write_url = table.unpack(tt)
        t:reset()
        t.env.c_log_active_ch = nil
        t.env.c_log_bool_ch_start = in_c_log_bool_ch_start
        t.env.c_log_bool_ch_limit = in_c_log_bool_ch_limit
        t.env.c_log_number_ch_start = in_c_log_number_start
        t.env.c_log_number_ch_limit = in_c_log_number_limit
        t.env.property._number["Port"] = 52149
        t.env.property._text["Title"] = "title"
        for k, v in pairs(in_property_text) do
            t.env.property._text[k] = v
        end
        t.fn()

        t.env.logOnTick()
        t.env.async:assertCall(1, 52149, "/time")

        t.env.httpReply(52149, "/time", "SVCOK20060102150405")
        t.env.async:assertCall(2, 52149, want_write_url)
    end
end

function g_test_tbl.testLogTick(t)
    t:reset()
    t.env.c_log_active_ch = 32
    t.env.c_log_bool_ch_start = 1
    t.env.c_log_bool_ch_limit = 0
    t.env.c_log_number_ch_start = 1
    t.env.c_log_number_ch_limit = 0
    t.env.input._bool[32] = true
    t.env.property._number["Port"] = 52149
    t.env.property._text["Title"] = "title"
    t.fn()

    t.env.logOnTick()
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=%23%0D%0A0%0D%0A")

    t.env.httpReply(52149, "/write?path=title%2Ftitle-20060102150405.csv&data=%23%0D%0A0%0D%0A", "SVCOK")
    t.env.async:assertCallCount(2)

    t.env.logOnTick()
    t.env.async:assertCall(3, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=1%0D%0A")

    t.env.httpReply(52149, "/write?path=title%2Ftitle-20060102150405.csv&data=1%0D%0A", "SVCOK")
    t.env.async:assertCallCount(3)

    t.env.logOnTick()
    t.env.async:assertCall(4, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=2%0D%0A")

    t.env.httpReply(52149, "/write?path=title%2Ftitle-20060102150405.csv&data=2%0D%0A", "SVCOK")
    t.env.async:assertCallCount(4)

    t.env.input._bool[32] = false
    t.env.logOnTick()
    t.env.async:assertCallCount(4)

    t.env.input._bool[32] = true
    t.env.logOnTick()
    t.env.async:assertCall(5, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(6, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=%23%0D%0A0%0D%0A")
end

function g_test_tbl.testLogRecord(t)
    local tests = {
        {
            1, 0, 1, 0,
            {},
            {},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%0D%0A",
        },
        {
            1, 1, 1, 0,
            {[1] = false},
            {},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2CFALSE%0D%0A",
        },
        {
            1, 1, 1, 0,
            {[1] = true},
            {},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2CTRUE%0D%0A",
        },
        {
            2, 3, 1, 0,
            {[2] = false, [3] = true},
            {},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2CFALSE%2CTRUE%0D%0A",
        },
        {
            1, 0, 1, 1,
            {},
            {[1] = 10},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2C10%0D%0A",
        },
        {
            1, 0, 1, 1,
            {},
            {[1] = 0.1},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2C0.1%0D%0A",
        },
        {
            1, 0, 1, 1,
            {},
            {[1] = 1000000000000000},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2C1E%2B15%0D%0A",
        },
        {
            1, 0, 2, 3,
            {},
            {[2] = 4, [3] = 5},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2C4%2C5%0D%0A",
        },
        {
            1, 1, 1, 1,
            {[1] = true},
            {[1] = 2},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2CTRUE%2C2%0D%0A",
        },
        {
            3, 4, 1, 2,
            {[3] = false, [4] = true},
            {[1] = 5, [2] = 6},
            "/write?path=title%2Ftitle-20060102150405.csv&data=1%2CFALSE%2CTRUE%2C5%2C6%0D%0A",
        },
    }

    for i, tt in ipairs(tests) do
        local in_c_log_bool_ch_start, in_c_log_bool_ch_limit, in_c_log_number_start, in_c_log_number_limit, in_input_bool, in_input_number, want_write_url = table.unpack(tt)
        t:reset()
        t.env.c_log_active_ch = nil
        t.env.c_log_bool_ch_start = in_c_log_bool_ch_start
        t.env.c_log_bool_ch_limit = in_c_log_bool_ch_limit
        t.env.c_log_number_ch_start = in_c_log_number_start
        t.env.c_log_number_ch_limit = in_c_log_number_limit
        t.env.property._number["Port"] = 52149
        t.env.property._text["Title"] = "title"
        t.fn()

        t.env.logOnTick()
        t.env.async:assertCall(1, 52149, "/time")

        t.env.httpReply(52149, "/time", "SVCOK20060102150405")
        t.env.async:assertCallCount(2)

        t.env.httpReply(52149, t.env.async._url, "SVCOK")
        t.env.async:assertCallCount(2)

        for k, v in pairs(in_input_bool) do
            t.env.input._bool[k] = v
        end
        for k, v in pairs(in_input_number) do
            t.env.input._number[k] = v
        end
        t.env.logOnTick()
        t.env.async:assertCall(3, 52149, want_write_url)
    end
end

function g_test_tbl.testSendRequestStateInit(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")
end

function g_test_tbl.testSendRequestStateCancel(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(2, 52150, "/time")

    t.env.httpReply(52150, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(3, 52150, "/write?path=title2%2Ftitle2-20060102150405.csv&data=data2")
end

function g_test_tbl.testSendRequestStateError(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title/", "data")
    t.env.async:assertCallCount(0)

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testSendRequestStateTime(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendRequest(52150, "title2", "data2")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1data2")
end

function g_test_tbl.testSendRequestStateWrite(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.sendRequest(52150, "title2", "data2")
    t.env.sendRequest(52151, "title3", "data3")

    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "SVCOK")
    t.env.async:assertCall(3, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data2data3")
end

function g_test_tbl.testSendRequestErrorTitleEmpty(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "", "data")
    t.env.async:assertCallCount(0)

    -- confirm send error
    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testSendRequestErrorTitleSlash(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title/", "data")
    t.env.async:assertCallCount(0)

    -- confirm send error
    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testSendRequestIgnorePortTitle(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendRequest(nil, nil, "data2")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1data2")
end

function g_test_tbl.testSendCancelStateInit(t)
    t:reset()
    t.fn()

    t.env.sendCancel()

    -- confirm send inactive
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(1, 52150, "/time")
end

function g_test_tbl.testSendCancelStateCancel(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()
    t.env.sendCancel()

    -- confirm client cancel
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    -- confirm send inactive
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(2, 52150, "/time")

    -- confirm send reset
    t.env.httpReply(52150, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(3, 52150, "/write?path=title2%2Ftitle2-20060102150405.csv&data=data2")
end

function g_test_tbl.testSendCancelStateError(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title/", "data1")
    t.env.async:assertCallCount(0)

    t.env.sendCancel()

    -- confirm send inactive
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(1, 52150, "/time")

    -- confirm send reset
    t.env.httpReply(52150, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52150, "/write?path=title2%2Ftitle2-20060102150405.csv&data=data2")
end

function g_test_tbl.testSendCancelStateTime(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()

    -- confirm client cancel
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    -- confirm send inactive
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(2, 52150, "/time")

    -- confirm send reset
    t.env.httpReply(52150, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(3, 52150, "/write?path=title2%2Ftitle2-20060102150405.csv&data=data2")
end

function g_test_tbl.testSendCancelStateWrite(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(2)

    t.env.sendCancel()

    -- confirm client cancel
    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "SVCOK")
    t.env.async:assertCallCount(2)

    -- confirm send inactive
    t.env.sendRequest(52151, "title3", "data3")
    t.env.async:assertCall(3, 52151, "/time")

    -- confirm send reset
    t.env.httpReply(52151, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(4, 52151, "/write?path=title3%2Ftitle3-20060102150405.csv&data=data3")
end

function g_test_tbl.testSendOnTickStateInit(t)
    t:reset()
    t.fn()

    t.env.sendOnTick()
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testSendOnTickStateCancel(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    t.env.sendOnTick()
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendOnTickStateError(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "timeout")
    t.env.async:assertCallCount(1)

    t.env.sendOnTick()
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendOnTickStateTime(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()
    t.env.sendRequest(52150, "title2", "data2")
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    t.env.sendOnTick()
    t.env.async:assertCall(2, 52150, "/time")
end

function g_test_tbl.testSendTimeRequestPend(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCall(1, 52149, "/time")

    -- confirm accept response
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=data")
end

function g_test_tbl.testSendTimeRequestBusy(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()
    t.env.sendRequest(52150, "title2", "data2")
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    -- confirm retry request
    t.env.sendOnTick()
    t.env.async:assertCall(2, 52150, "/time")
end

function g_test_tbl.testSendTimeResponseDoneOK(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=data")
end

function g_test_tbl.testSendTimeResponseDoneSlash(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK2006010215040/")
    t.env.async:assertCallCount(1)

    -- confirm send error
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendTimeResponseDoneError(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "timeout")
    t.env.async:assertCallCount(1)

    -- confirm send error
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendTimeResponseCancel(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendCancel()
    t.env.async:assertCallCount(1)

    -- confirm send reset
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(2, 52150, "/time")
end

function g_test_tbl.testSendTimeResponseTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 1

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.sendOnTick()
    t.env.async:assertCallCount(1)

    -- confirm client timeout
    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)

    -- confirm send error
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendWriteRequestPend(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title%2Ftitle-20060102150405.csv&data=data")

    -- confirm buf clear
    t.env.httpReply(52149, "/write?path=title%2Ftitle-20060102150405.csv&data=data", "SVCOK")
    t.env.async:assertCallCount(2)
end

function g_test_tbl.testSendWriteRequestSize(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 53

    t.env.sendRequest(52149, "title", "data")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testSendWriteResponseDoneOK(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "SVCOK")
    t.env.async:assertCallCount(2)

    -- confirm send active
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(3, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data2")
end

function g_test_tbl.testSendWriteResponseDoneError(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "timeout")
    t.env.async:assertCallCount(2)

    -- confirm send error
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(2)
end

function g_test_tbl.testSendWriteResponseCancel(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.sendCancel()
    t.env.async:assertCallCount(2)

    -- confirm send reset
    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "SVCOK")
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCall(3, 52150, "/time")
end

function g_test_tbl.testSendWriteResponseTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 1

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")

    t.env.sendOnTick()
    t.env.async:assertCallCount(2)

    -- confirm send error
    t.env.httpReply(52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1", "SVCOK")
    t.env.sendRequest(52150, "title2", "data2")
    t.env.async:assertCallCount(2)
end

function g_test_tbl.testSendEventTime(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")
end

function g_test_tbl.testSendEventWrite(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "data1")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCall(2, 52149, "/write?path=title1%2Ftitle1-20060102150405.csv&data=data1")
end

function g_test_tbl.testSendEventWriteEmpty(t)
    t:reset()
    t.fn()

    t.env.sendRequest(52149, "title1", "")
    t.env.async:assertCall(1, 52149, "/time")

    t.env.httpReply(52149, "/time", "SVCOK20060102150405")
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testClientSizeNormal(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback("callback")

    local status = t.env.clientHttpGet("ctx", 52149, "/ur", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
    callback:assertWait()
    t.env.async:assertCall(1, 52149, "/ur")
end

function g_test_tbl.testClientSizeError(t)
    t:reset()
    t.fn()

    t.env.c_client_maxlen = 3
    local callback = buildMockClientCallback("callback")

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_size, status)
    callback:assertWait()
    t.env.async:assertCallCount(0)
end

function g_test_tbl.testClientBusyAfterInit(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
    callback:assertWait()
    t.env.async:assertCall(1, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterGet(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_busy, status)
    callback:assertWait()
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testClientBusyAfterCancel(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpCancel()

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_busy, status)
    callback:assertWait()
    t.env.async:assertCallCount(1)
end

function g_test_tbl.testClientBusyAfterTimeout(t)
    t:reset()
    t.fn()

    t.env.c_client_timeout = 0
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientOnTick()

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
    callback:assertWait()
    t.env.async:assertCall(2, 52149, "/url")
end

function g_test_tbl.testClientBusyAfterReply(t)
    t:reset()
    t.fn()
    local callback = buildMockClientCallback("callback")

    t.env.clientHttpGet("ctx", 52149, "/url", function() end)
    t.env.clientHttpReply(52149, "/url", "resp")

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
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

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
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

    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_pend, status)
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
    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_busy, status)
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
    callback:assertCall("ctx", t.env.c_client_status_done, "resp")

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
    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_busy, status)
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
    local status = t.env.clientHttpGet("ctx", 52149, "/url", callback.fn)
    assertEqual("status", t.env.c_client_status_busy, status)
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
            self.env.input = buildMockInput()
            self.env.output = buildMockOutput()
            self.env.property = buildMockProperty()
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

function buildMockInput()
    local input = {
        _bool = {},
        _number = {},
    }

    input.getNumber = function(index)
        return input._number[index] or 0
    end

    input.getBool = function(index)
        return input._bool[index] or false
    end

    return input
end

function buildMockOutput()
    local output = {}

    output.setBool = function(index, value)
    end

    output.setNumber = function(index, value)
    end

    return output
end

function buildMockProperty()
    local property = {
        _number = {},
        _bool = {},
        _text = {},
    }

    property.getNumber = function(label)
        return property._number[label] or 0
    end

    property.getBool = function(label)
        return property._bool[label] or false
    end

    property.getText = function(label)
        return property._text[label] or ""
    end

    return property
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
