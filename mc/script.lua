function escapeQuery(s)
    s = tostring(s)

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

function encodeCSVRecord(record)
    -- RFC 4180
    local out = {}
    for i, s in ipairs(record) do
        out[i] = escapeCSVField(s)
    end
    return table.concat(out, ",") .. "\r\n"
end

function escapeCSVField(s)
    -- RFC 4180
    if string.match(s, "\r\n") ~= nil or string.match(s, '[",]') ~= nil then
        s = string.gsub(s, '"', '""')
        s = '"' .. s .. '"'
    end
    return s
end
