package rate_limit

const (
	RateLimitIP = "rate:limit:ip_"
)

const LimitLuaScript = `
local current = redis.call("get", KEYS[1])
if current and tonumber(current) >= tonumber(ARGV[1]) then
	return 0
else
	current = redis.call("incr", KEYS[1])
	if tonumber(current) == 1 then
		redis.call("expire", KEYS[1], ARGV[2])
	end
	return 1
end
`
