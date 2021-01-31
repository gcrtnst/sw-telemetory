function init()
	-- property --
	p_active = property.getBool('Active')
	p_port = property.getNumber('Port')
	p_title = property.getText('Title')
	p_numcol = property.getNumber('Number of Columns')

	p_label = {}
	for i = 1, p_numcol do
		p_label[i] = property.getText('Label ' .. string.format('%d', i))
	end

	-- global --
	g_active = false
	g_tick = 0

	g_header = {''}
	for i = 1, p_numcol do
		local label = p_label[i]
		label = encodeCSVField(label)
		table.insert(g_header, label)
	end
	g_header = table.concat(g_header, ',')

	g_client = buildClient()
	g_client['sender']['port'] = p_port
end

function onTick()
	-- active --
	local active = input.getBool(1)
	if g_active ~= (p_active or active) then
		g_active = p_active or active
		if g_active then
			g_client.new(p_title)
			g_client.write(g_header)
		end
	end

	-- write --
	if g_active then
		local buf = {string.format('%d', g_tick)}
		for i = 1, p_numcol do
			local num = input.getNumber(i)
			table.insert(buf, string.format('%G', num))
		end
		buf = table.concat(buf, ',')
		g_client.write(buf)
	end

	-- tick --
	if g_active then
		g_tick = g_tick + 1
	else
		g_tick = 0
	end
end


function httpReply(port, request, response)
	g_client['sender'].httpReply(port, request, response)
end

function encodeCSVField(s)
	if string.match(s, '[\n",]') ~= fail then
		s = string.gsub(s, '"', '""')
		s = '"' .. s .. '"'
	end
	return s
end

function buildClient()
	local client = {
		['sender'] = buildSender(),
	}

	function client.new(title)
		title = string.gsub(title, '\n', '')
		client['sender'].send('n' .. title .. '\n')
	end

	function client.write(s)
		s = string.gsub(s, '\n', '')
		client['sender'].send('w' .. s .. '\n')
	end

	return client
end

function buildSender()
	local sender = {
		['port'] = 58592,
		['_buf'] = {},
		['_buf_size'] = 0,
		['_sending'] = false,
	}

	function sender.send(s)
		local max_buf_size = 498

		if #s <= 0 or #s > max_buf_size then
			return
		end
		sender['_buf_size'] = sender['_buf_size'] + #s
		table.insert(sender['_buf'], s)
		while sender['_buf_size'] > max_buf_size do
			sender['_buf_size'] = sender['_buf_size'] - #sender['_buf'][1]
			table.remove(sender['_buf'], 1)
		end

		if not sender['_sending'] then
			sender._send()
		end
	end

	function sender.httpReply(port, request, reply)
		if port ~= sender['port'] then
			return
		end
		sender['_sending'] = false

		if sender['_buf_size'] > 0 then
			sender._send()
		end
	end

	function sender._send()
		local buf = table.concat(sender['_buf'])
		async.httpGet(sender['port'], buf)
		sender['_buf'] = {}
		sender['_buf_size'] = 0
		sender['_sending'] = true
	end

	return sender
end

init()
