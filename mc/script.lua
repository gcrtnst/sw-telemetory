function init()
    sendInit()
end

function onTick()
    sendOnTick()
end

function httpReply(port, req, resp)
    clientHttpReply(port, req, resp)
end

function sendInit()
    clientInit()

    c_send_ctx_time = 0x11
    c_send_ctx_write = 0x12
    sendInitVar()
end

function sendInitVar()
    g_send_error = false
    g_send_active = false
    g_send_port = nil
    g_send_title = nil
    g_send_buf = nil
    g_send_path = nil
end

function sendOnTick()
    clientOnTick()

    if g_send_active then
        sendEvent()
    end
end

function sendRequest(port, title, data)
    if g_send_error then
        return
    end

    if not g_send_active then
        if title == "" or string.match(title, "/") ~= nil then
            sendError()
            return
        end

        g_send_active = true
        g_send_port = port
        g_send_title = title
        g_send_buf = ""
    end

    g_send_buf = g_send_buf .. data
    sendEvent()
end

function sendCancel()
    clientHttpCancel()
    sendInitVar()
end

function sendEvent()
    if g_send_path == nil then
        local status = clientHttpGet(c_send_ctx_time, g_send_port, "/time", sendCallback)
        if status ~= c_client_status_pend and status ~= c_client_status_busy then
            sendError()
            return
        end
        return
    end

    if #g_send_buf > 0 then
        local req = "/write?path=" .. escapeQuery(g_send_path) .. "&data=" .. escapeQuery(g_send_buf)
        local status = clientHttpGet(c_send_ctx_write, g_send_port, req, sendCallback)
        if status == c_client_status_busy then
            return
        end
        if status ~= c_client_status_pend then
            sendError()
            return
        end
        g_send_buf = ""
        return
    end
end

function sendCallback(ctx, status, resp)
    if status ~= c_client_status_done or string.sub(resp, 1, 5) ~= "SVCOK" then
        sendError()
        return
    end

    if ctx == c_send_ctx_time then
        local time = string.sub(resp, 6)
        if string.match(time, "^%d%d%d%d%d%d%d%d%d%d%d%d%d%d$") == nil then
            sendError()
            return
        end
        g_send_path = string.format("%s/%s-%s.csv", g_send_title, g_send_title, time)
    end

    sendEvent()
end

function sendError()
    sendCancel()
    g_send_error = true
end

function clientInit()
    c_client_maxlen = 3840
    c_client_timeout = 600

    c_client_status_done = 0x00
    c_client_status_pend = 0x01
    c_client_status_size = 0x02
    c_client_status_busy = 0x03
    c_client_status_cancel = 0x04
    c_client_status_timeout = 0x05

    clientInitVar()
end

function clientInitVar()
    g_client_timeout = nil
    g_client_ctx = nil
    g_client_port = nil
    g_client_req = nil
    g_client_callback = nil
end

function clientOnTick()
    if g_client_timeout ~= nil then
        g_client_timeout = g_client_timeout - 1
        if g_client_timeout <= 0 then
            clientHttpFinish(c_client_status_timeout, nil)
            return
        end
    end
end

function clientHttpReply(port, req, resp)
    if g_client_timeout == nil or g_client_port ~= port or g_client_req ~= req then
        return
    end
    clientHttpFinish(c_client_status_done, resp)
end

function clientHttpGet(ctx, port, req, callback)
    if #req > c_client_maxlen then
        return c_client_status_size
    end
    if g_client_timeout ~= nil then
        return c_client_status_busy
    end

    g_client_timeout = c_client_timeout
    g_client_ctx = ctx
    g_client_port = port
    g_client_req = req
    g_client_callback = callback
    async.httpGet(port, req)
    return c_client_status_pend
end

function clientHttpCancel()
    if g_client_timeout == nil then
        return
    end

    local ctx = g_client_ctx
    local callback = g_client_callback
    g_client_ctx = nil
    g_client_callback = function() end

    callback(ctx, c_client_status_cancel, nil)
end

function clientHttpFinish(status, resp)
    local ctx = g_client_ctx
    local callback = g_client_callback
    clientInitVar()

    callback(ctx, status, resp)
end

function encodeCSVRecord(record)
    if type(record) ~= "table" then
        return nil
    end

    -- RFC 4180
    local out = {}
    for i, s in ipairs(record) do
        local o = encodeCSVField(s)
        if o == nil then
            return nil
        end
        out[i] = o
    end
    return table.concat(out, ",") .. "\r\n"
end

function encodeCSVField(s)
    if type(s) ~= "string" then
        return nil
    end

    -- RFC 4180
    if string.match(s, "\r\n") ~= nil or string.match(s, '[",]') ~= nil then
        s = string.gsub(s, '"', '""')
        s = '"' .. s .. '"'
    end
    return s
end

function escapeQuery(s)
    if type(s) ~= "string" then
        return nil
    end

    local out = {}
    for i = 1, #s do
        local c = string.byte(s, i)
        local o
        if (
            c == 0x2D or                    -- -
            c == 0x2E or                    -- .
            (0x30 <= c and c <= 0x39) or    -- 0..9
            (0x41 <= c and c <= 0x5A) or    -- A..Z
            c == 0x5F or                    -- _
            (0x61 <= c and c <= 0x7A) or    -- a..z
            c == 0x7E                       -- ~
        ) then
            o = string.char(c)
        elseif c == 0x20 then   -- space
            o = "+"
        else
            o = string.format("%%%02X", c)
        end
        table.insert(out, o)
    end
    return table.concat(out)
end

init()
