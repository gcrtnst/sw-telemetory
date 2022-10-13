function init()
    clientInit()
end

function onTick()
    clientOnTick()
end

function clientInit()
    c_client_maxlen = 3840
    c_client_timeout = 600

    c_client_status_done = "done"
    c_client_status_pend = "pend"
    c_client_status_size = "size"
    c_client_status_busy = "busy"
    c_client_status_cancel = "cancel"
    c_client_status_timeout = "timeout"

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
