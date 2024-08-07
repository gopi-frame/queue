package redis

// LuaScript is a Lua script set for redis queue
type LuaScript struct{}

// Enqueue returns the Lua script for pushing a job
//
// # Parameters:
//
//	KEYS[1] = index sorted set key
//	KEYS[2] = hash table key
//	ARGV[1] = score
//	ARGV[2] = job id
//	ARGV[3] = job hash
func (LuaScript) Enqueue() string {
	return `
-- add job id into index sorted set
redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2])

-- add job hash into hash table
redis.call('HSET', KEYS[2], ARGV[2], ARGV[3])
`
}

// Dequeue returns the Lua script for popping a job
//
// Parameters:
//
//	KEYS[1] = index sorted set key
//	KEYS[2] = hash table key
//	ARGV[1] = number of jobs to pop
//	ARGV[2] = current time
func (LuaScript) Dequeue() string {
	return `
-- pop job id from index sorted set
local result = redis.call('ZPOPMIN', KEYS[1], ARGV[1])

-- return nil if no job is found
if next(result) == nil then
	return nil
end

-- return nil if availability time is not reached
if tonumber(result[2]) > tonumber(ARGV[2]) then
	-- push job id back to index sorted set
	redis.call('ZADD', KEYS[1], result[2], result[1])
	return nil
end



local id = result[1]

-- get job hash from hash table
local job = redis.call('HGET', KEYS[2], id)

-- remove job hash from hash table
redis.call('HDEL', KEYS[2], id)

return job
`
}

// Remove returns the Lua script for removing a job
//
// Parameters:
//
//	KEYS[1] = index sorted set key
//	KEYS[2] = hash table key
//	ARGV[1] = job id
func (LuaScript) Remove() string {
	return `
-- remove job id from index sorted set
redis.call('ZREM', KEYS[1], ARGV[1])

-- remove job hash from hash table
redis.call('HDEL', KEYS[2], ARGV[1])
`
}
