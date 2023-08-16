package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func (cli *Client) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Pipelined(ctx, fn)
}

func (cli *Client) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.TxPipelined(ctx, fn)
}

func (cli *Client) Command() *redis.CommandsInfoCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Command(ctx)
}

func (cli *Client) ClientGetName() *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClientGetName(ctx)
}

func (cli *Client) Echo(message interface{}) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Echo(ctx, message)
}

func (cli *Client) Ping() *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Ping(ctx)
}

func (cli *Client) Quit() *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Quit(ctx)
}

func (cli *Client) Del(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Del(ctx, keys...)
}

func (cli *Client) Unlink(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Unlink(ctx, keys...)
}

func (cli *Client) Dump(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Dump(ctx, key)
}

func (cli *Client) Exists(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Exists(ctx, keys...)
}

func (cli *Client) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Expire(ctx, key, expiration)
}

func (cli *Client) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ExpireAt(ctx, key, tm)
}

func (cli *Client) ExpireNX(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ExpireNX(ctx, key, expiration)
}

func (cli *Client) ExpireXX(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ExpireXX(ctx, key, expiration)
}

func (cli *Client) ExpireGT(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ExpireGT(ctx, key, expiration)
}

func (cli *Client) ExpireLT(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ExpireLT(ctx, key, expiration)
}

func (cli *Client) Keys(pattern string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Keys(ctx, pattern)
}

func (cli *Client) Migrate(host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Migrate(ctx, host, port, key, db, timeout)
}

func (cli *Client) Move(key string, db int) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Move(ctx, key, db)
}

func (cli *Client) ObjectRefCount(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ObjectRefCount(ctx, key)
}

func (cli *Client) ObjectEncoding(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ObjectEncoding(ctx, key)
}

func (cli *Client) ObjectIdleTime(key string) *redis.DurationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ObjectIdleTime(ctx, key)
}

func (cli *Client) Persist(key string) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Persist(ctx, key)
}

func (cli *Client) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PExpire(ctx, key, expiration)
}

func (cli *Client) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PExpireAt(ctx, key, tm)
}

func (cli *Client) PTTL(key string) *redis.DurationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PTTL(ctx, key)
}

func (cli *Client) Rename(key, newkey string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Rename(ctx, key, newkey)
}

func (cli *Client) RenameNX(key, newkey string) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RenameNX(ctx, key, newkey)
}

func (cli *Client) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Restore(ctx, key, ttl, value)
}

func (cli *Client) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RestoreReplace(ctx, key, ttl, value)
}

func (cli *Client) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Sort(ctx, key, sort)
}

func (cli *Client) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SortStore(ctx, key, store, sort)
}

func (cli *Client) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SortInterfaces(ctx, key, sort)
}

func (cli *Client) Touch(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Touch(ctx, keys...)
}

func (cli *Client) TTL(key string) *redis.DurationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.TTL(ctx, key)
}

func (cli *Client) Type(key string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Type(ctx, key)
}

func (cli *Client) Append(key, value string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Append(ctx, key, value)
}

func (cli *Client) Decr(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Decr(ctx, key)
}

func (cli *Client) DecrBy(key string, decrement int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.DecrBy(ctx, key, decrement)
}

func (cli *Client) Get(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Get(ctx, key)
}

func (cli *Client) GetRange(key string, start, end int64) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GetRange(ctx, key, start, end)
}

func (cli *Client) GetSet(key string, value interface{}) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GetSet(ctx, key, value)
}

func (cli *Client) GetEx(key string, expiration time.Duration) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GetEx(ctx, key, expiration)
}

func (cli *Client) GetDel(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GetDel(ctx, key)
}

func (cli *Client) Incr(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Incr(ctx, key)
}

func (cli *Client) IncrBy(key string, value int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.IncrBy(ctx, key, value)
}

func (cli *Client) IncrByFloat(key string, value float64) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.IncrByFloat(ctx, key, value)
}

func (cli *Client) MGet(keys ...string) *redis.SliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.MGet(ctx, keys...)
}

func (cli *Client) MSet(values ...interface{}) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.MSet(ctx, values...)
}

func (cli *Client) MSetNX(values ...interface{}) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.MSetNX(ctx, values...)
}

func (cli *Client) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Set(ctx, key, value, expiration)
}

func (cli *Client) SetArgs(key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetArgs(ctx, key, value, a)
}

func (cli *Client) SetEX(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetEX(ctx, key, value, expiration)
}

func (cli *Client) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetNX(ctx, key, value, expiration)
}

func (cli *Client) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetXX(ctx, key, value, expiration)
}

func (cli *Client) SetRange(key string, offset int64, value string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetRange(ctx, key, offset, value)
}

func (cli *Client) StrLen(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.StrLen(ctx, key)
}

func (cli *Client) Copy(sourceKey string, destKey string, db int, replace bool) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Copy(ctx, sourceKey, destKey, db, replace)
}

func (cli *Client) GetBit(key string, offset int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GetBit(ctx, key, offset)
}

func (cli *Client) SetBit(key string, offset int64, value int) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SetBit(ctx, key, offset, value)
}

func (cli *Client) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitCount(ctx, key, bitCount)
}

func (cli *Client) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitOpAnd(ctx, destKey, keys...)
}

func (cli *Client) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitOpOr(ctx, destKey, keys...)
}

func (cli *Client) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitOpXor(ctx, destKey, keys...)
}

func (cli *Client) BitOpNot(destKey string, key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitOpNot(ctx, destKey, key)
}

func (cli *Client) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitPos(ctx, key, bit, pos...)
}

func (cli *Client) BitField(key string, args ...interface{}) *redis.IntSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BitField(ctx, key, args)
}

func (cli *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return cli.Client.Scan(ctx, cursor, match, count)
}

func (cli *Client) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return cli.Client.ScanType(ctx, cursor, match, count, keyType)
}

func (cli *Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return cli.Client.SScan(ctx, key, cursor, match, count)
}

func (cli *Client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return cli.Client.HScan(ctx, key, cursor, match, count)
}

func (cli *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return cli.Client.ZScan(ctx, key, cursor, match, count)
}

func (cli *Client) HDel(key string, fields ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HDel(ctx, key, fields...)
}

func (cli *Client) HExists(key, field string) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HExists(ctx, key, field)
}

func (cli *Client) HGet(key, field string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HGet(ctx, key, field)
}

func (cli *Client) HGetAll(key string) *redis.StringStringMapCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HGetAll(ctx, key)
}

func (cli *Client) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HIncrBy(ctx, key, field, incr)
}

func (cli *Client) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HIncrByFloat(ctx, key, field, incr)
}

func (cli *Client) HKeys(key string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HKeys(ctx, key)
}

func (cli *Client) HLen(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HLen(ctx, key)
}

func (cli *Client) HMGet(key string, fields ...string) *redis.SliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HMGet(ctx, key, fields...)
}

func (cli *Client) HSet(key string, values ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HSet(ctx, key, values...)
}

func (cli *Client) HMSet(key string, values ...interface{}) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HMSet(ctx, key, values...)
}

func (cli *Client) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HSetNX(ctx, key, field, value)
}

func (cli *Client) HVals(key string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HVals(ctx, key)
}

func (cli *Client) HRandField(key string, count int, withValues bool) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.HRandField(ctx, key, count, withValues)
}

func (cli *Client) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BLPop(ctx, timeout, keys...)
}

func (cli *Client) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BRPop(ctx, timeout, keys...)
}

func (cli *Client) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BRPopLPush(ctx, source, destination, timeout)
}

func (cli *Client) LIndex(key string, index int64) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LIndex(ctx, key, index)
}

func (cli *Client) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LInsert(ctx, key, op, pivot, value)
}

func (cli *Client) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LInsertBefore(ctx, key, pivot, value)
}

func (cli *Client) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LInsertAfter(ctx, key, pivot, value)
}

func (cli *Client) LLen(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LLen(ctx, key)
}

func (cli *Client) LPop(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPop(ctx, key)
}

func (cli *Client) LPopCount(key string, count int) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPopCount(ctx, key, count)
}

func (cli *Client) LPos(key string, value string, args redis.LPosArgs) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPos(ctx, key, value, args)
}

func (cli *Client) LPosCount(key string, value string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPosCount(ctx, key, value, count, args)
}

func (cli *Client) LPush(key string, values ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPush(ctx, key, values...)
}

func (cli *Client) LPushX(key string, values ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LPushX(ctx, key, values...)
}

func (cli *Client) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LRange(ctx, key, start, stop)
}

func (cli *Client) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LRem(ctx, key, count, value)
}

func (cli *Client) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LSet(ctx, key, index, value)
}

func (cli *Client) LTrim(key string, start, stop int64) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LTrim(ctx, key, start, stop)
}

func (cli *Client) RPop(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RPop(ctx, key)
}

func (cli *Client) RPopCount(key string, count int) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RPopCount(ctx, key, count)
}

func (cli *Client) RPopLPush(source, destination string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RPopLPush(ctx, source, destination)
}

func (cli *Client) RPush(key string, values ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RPush(ctx, key, values...)
}

func (cli *Client) RPushX(key string, values ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.RPushX(ctx, key, values...)
}

func (cli *Client) LMove(source, destination, srcpos, destpos string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.LMove(ctx, source, destination, srcpos, destpos)
}

func (cli *Client) BLMove(source, destination, srcpos, destpos string, timeout time.Duration) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BLMove(ctx, source, destination, srcpos, destpos, timeout)
}

func (cli *Client) SAdd(key string, members ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SAdd(ctx, key, members...)
}

func (cli *Client) SCard(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SCard(ctx, key)
}

func (cli *Client) SDiff(keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SDiff(ctx, keys...)
}

func (cli *Client) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SDiffStore(ctx, destination, keys...)
}

func (cli *Client) SInter(keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SInter(ctx, keys...)
}

func (cli *Client) SInterStore(destination string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SInterStore(ctx, destination, keys...)
}

func (cli *Client) SIsMember(key string, member interface{}) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SIsMember(ctx, key, member)
}

func (cli *Client) SMIsMember(key string, members ...interface{}) *redis.BoolSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SMIsMember(ctx, key, members...)
}

func (cli *Client) SMembers(key string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SMembers(ctx, key)
}

func (cli *Client) SMembersMap(key string) *redis.StringStructMapCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SMembersMap(ctx, key)
}

func (cli *Client) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SMove(ctx, source, destination, member)
}

func (cli *Client) SPop(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SPop(ctx, key)
}

func (cli *Client) SPopN(key string, count int64) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SPopN(ctx, key, count)
}

func (cli *Client) SRandMember(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SRandMember(ctx, key)
}

func (cli *Client) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SRandMemberN(ctx, key, count)
}

func (cli *Client) SRem(key string, members ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SRem(ctx, key, members)
}

func (cli *Client) SUnion(keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SUnion(ctx, keys...)
}

func (cli *Client) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SUnionStore(ctx, destination, keys...)
}

func (cli *Client) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XAdd(ctx, a)
}

func (cli *Client) XDel(stream string, ids ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XDel(ctx, stream, ids...)
}

func (cli *Client) XLen(stream string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XLen(ctx, stream)
}

func (cli *Client) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XRange(ctx, stream, start, stop)
}

func (cli *Client) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XRangeN(ctx, stream, start, stop, count)
}

func (cli *Client) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XRevRange(ctx, stream, start, stop)
}

func (cli *Client) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XRevRangeN(ctx, stream, start, stop, count)
}

func (cli *Client) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XRead(ctx, a)
}

func (cli *Client) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XReadStreams(ctx, streams...)
}

func (cli *Client) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupCreate(ctx, stream, group, start)
}

func (cli *Client) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupCreateMkStream(ctx, stream, group, start)
}

func (cli *Client) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupSetID(ctx, stream, group, start)
}

func (cli *Client) XGroupDestroy(stream, group string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupDestroy(ctx, stream, group)
}

func (cli *Client) XGroupCreateConsumer(stream, group, consumer string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupCreateConsumer(ctx, stream, group, consumer)
}

func (cli *Client) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XGroupDelConsumer(ctx, stream, group, consumer)
}

func (cli *Client) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XReadGroup(ctx, a)
}

func (cli *Client) XAck(stream, group string, ids ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XAck(ctx, stream, group, ids...)
}

func (cli *Client) XPending(stream, group string) *redis.XPendingCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XPending(ctx, stream, group)
}

func (cli *Client) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XPendingExt(ctx, a)
}

func (cli *Client) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XClaim(ctx, a)
}

func (cli *Client) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XClaimJustID(ctx, a)
}

func (cli *Client) XAutoClaim(a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XAutoClaim(ctx, a)
}

func (cli *Client) XAutoClaimJustID(a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XAutoClaimJustID(ctx, a)
}

func (cli *Client) XTrim(key string, maxLen int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrim(ctx, key, maxLen)
}

func (cli *Client) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrimApprox(ctx, key, maxLen)
}

func (cli *Client) XTrimMaxLen(key string, maxLen int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrimMaxLen(ctx, key, maxLen)
}

func (cli *Client) XTrimMaxLenApprox(key string, maxLen, limit int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrimMaxLenApprox(ctx, key, maxLen, limit)
}

func (cli *Client) XTrimMinID(key string, minID string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrimMinID(ctx, key, minID)
}

func (cli *Client) XTrimMinIDApprox(key string, minID string, limit int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XTrimMinIDApprox(ctx, key, minID, limit)
}

func (cli *Client) XInfoGroups(key string) *redis.XInfoGroupsCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XInfoGroups(ctx, key)
}

func (cli *Client) XInfoStream(key string) *redis.XInfoStreamCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XInfoStream(ctx, key)
}

func (cli *Client) XInfoStreamFull(key string, count int) *redis.XInfoStreamFullCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XInfoStreamFull(ctx, key, count)
}

func (cli *Client) XInfoConsumers(key string, group string) *redis.XInfoConsumersCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.XInfoConsumers(ctx, key, group)
}

func (cli *Client) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BZPopMax(ctx, timeout, keys...)
}

func (cli *Client) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.BZPopMin(ctx, timeout, keys...)
}

func (cli *Client) ZAdd(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAdd(ctx, key, members...)
}

func (cli *Client) ZAddNX(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddNX(ctx, key, members...)
}

func (cli *Client) ZAddXX(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddXX(ctx, key, members...)
}

func (cli *Client) ZAddCh(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddCh(ctx, key, members...)
}

func (cli *Client) ZAddNXCh(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddNXCh(ctx, key, members...)
}

func (cli *Client) ZAddXXCh(key string, members ...*redis.Z) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddXXCh(ctx, key, members...)
}

func (cli *Client) ZAddArgs(key string, args redis.ZAddArgs) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddArgs(ctx, key, args)
}

func (cli *Client) ZAddArgsIncr(key string, args redis.ZAddArgs) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZAddArgsIncr(ctx, key, args)
}

func (cli *Client) ZIncr(key string, member *redis.Z) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZIncr(ctx, key, member)
}

func (cli *Client) ZIncrNX(key string, member *redis.Z) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZIncrNX(ctx, key, member)
}

func (cli *Client) ZIncrXX(key string, member *redis.Z) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZIncrXX(ctx, key, member)
}

func (cli *Client) ZCard(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZCard(ctx, key)
}

func (cli *Client) ZCount(key, min, max string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZCount(ctx, key, min, max)
}

func (cli *Client) ZLexCount(key, min, max string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZLexCount(ctx, key, min, max)
}

func (cli *Client) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZIncrBy(ctx, key, increment, member)
}

func (cli *Client) ZInter(store *redis.ZStore) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZInter(ctx, store)
}

func (cli *Client) ZInterWithScores(store *redis.ZStore) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZInterWithScores(ctx, store)
}

func (cli *Client) ZInterStore(destination string, store *redis.ZStore) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZInterStore(ctx, destination, store)
}

func (cli *Client) ZMScore(key string, members ...string) *redis.FloatSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZMScore(ctx, key, members...)
}

func (cli *Client) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZPopMax(ctx, key, count...)
}

func (cli *Client) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZPopMin(ctx, key, count...)
}

func (cli *Client) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRange(ctx, key, start, stop)
}

func (cli *Client) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeWithScores(ctx, key, start, stop)
}

func (cli *Client) ZRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeByScore(ctx, key, opt)
}

func (cli *Client) ZRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeByLex(ctx, key, opt)
}

func (cli *Client) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeByScoreWithScores(ctx, key, opt)
}

func (cli *Client) ZRangeArgs(z redis.ZRangeArgs) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeArgs(ctx, z)
}

func (cli *Client) ZRangeArgsWithScores(z redis.ZRangeArgs) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeArgsWithScores(ctx, z)
}

func (cli *Client) ZRangeStore(dst string, z redis.ZRangeArgs) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRangeStore(ctx, dst, z)
}

func (cli *Client) ZRank(key, member string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRank(ctx, key, member)
}

func (cli *Client) ZRem(key string, members ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRem(ctx, key, members...)
}

func (cli *Client) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRemRangeByRank(ctx, key, start, stop)
}

func (cli *Client) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRemRangeByScore(ctx, key, min, max)
}

func (cli *Client) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRemRangeByLex(ctx, key, min, max)
}

func (cli *Client) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRange(ctx, key, start, stop)
}

func (cli *Client) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRangeWithScores(ctx, key, start, stop)
}

func (cli *Client) ZRevRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRangeByScore(ctx, key, opt)
}

func (cli *Client) ZRevRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRangeByLex(ctx, key, opt)
}

func (cli *Client) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (cli *Client) ZRevRank(key, member string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRevRank(ctx, key, member)
}

func (cli *Client) ZScore(key, member string) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZScore(ctx, key, member)
}

func (cli *Client) ZUnionStore(dest string, store *redis.ZStore) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZUnionStore(ctx, dest, store)
}

func (cli *Client) ZUnion(store redis.ZStore) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZUnion(ctx, store)
}

func (cli *Client) ZUnionWithScores(store redis.ZStore) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZUnionWithScores(ctx, store)
}

func (cli *Client) ZRandMember(key string, count int, withScores bool) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZRandMember(ctx, key, count, withScores)
}

func (cli *Client) ZDiff(keys ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZDiff(ctx, keys...)
}

func (cli *Client) ZDiffWithScores(keys ...string) *redis.ZSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZDiffWithScores(ctx, keys...)
}

func (cli *Client) ZDiffStore(destination string, keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ZDiffStore(ctx, destination, keys...)
}

func (cli *Client) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PFAdd(ctx, key, els...)
}

func (cli *Client) PFCount(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PFCount(ctx, keys...)
}

func (cli *Client) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PFMerge(ctx, dest, keys...)
}

func (cli *Client) ClientKill(ipPort string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClientKill(ctx, ipPort)
}

func (cli *Client) ClientKillByFilter(keys ...string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClientKillByFilter(ctx, keys...)
}

func (cli *Client) ClientPause(dur time.Duration) *redis.BoolCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClientPause(ctx, dur)
}

func (cli *Client) ConfigGet(parameter string) *redis.SliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ConfigGet(ctx, parameter)
}

func (cli *Client) ConfigSet(parameter, value string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ConfigSet(ctx, parameter, value)
}

func (cli *Client) Info(section ...string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Info(ctx, section...)
}

func (cli *Client) SlaveOf(host, port string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.SlaveOf(ctx, host, port)
}

func (cli *Client) DebugObject(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.DebugObject(ctx, key)
}

func (cli *Client) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.MemoryUsage(ctx, key, samples...)
}

func (cli *Client) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Eval(ctx, script, keys, args...)
}

func (cli *Client) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.EvalSha(ctx, sha1, keys, args...)
}

func (cli *Client) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ScriptExists(ctx, hashes...)
}

func (cli *Client) ScriptLoad(script string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ScriptLoad(ctx, script)
}

func (cli *Client) Publish(channel string, message interface{}) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.Publish(ctx, channel, message)
}

func (cli *Client) PubSubChannels(pattern string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PubSubChannels(ctx, pattern)
}

func (cli *Client) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.PubSubNumSub(ctx, channels...)
}

func (cli *Client) ClusterMeet(host, port string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterMeet(ctx, host, port)
}

func (cli *Client) ClusterForget(nodeID string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterForget(ctx, nodeID)
}

func (cli *Client) ClusterReplicate(nodeID string) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterReplicate(ctx, nodeID)
}

func (cli *Client) ClusterKeySlot(key string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterKeySlot(ctx, key)
}

func (cli *Client) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterGetKeysInSlot(ctx, slot, count)
}

func (cli *Client) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterCountFailureReports(ctx, nodeID)
}

func (cli *Client) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterCountKeysInSlot(ctx, slot)
}

func (cli *Client) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterDelSlots(ctx, slots...)
}

func (cli *Client) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterDelSlotsRange(ctx, min, max)
}

func (cli *Client) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterSlaves(ctx, nodeID)
}

func (cli *Client) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterAddSlots(ctx, slots...)
}

func (cli *Client) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.ClusterAddSlotsRange(ctx, min, max)
}

func (cli *Client) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoAdd(ctx, key, geoLocation...)
}

func (cli *Client) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoPos(ctx, key, members...)
}

func (cli *Client) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoRadius(ctx, key, longitude, latitude, query)
}

func (cli *Client) GeoRadiusStore(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoRadiusStore(ctx, key, longitude, latitude, query)
}

func (cli *Client) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoRadiusByMember(ctx, key, member, query)
}

func (cli *Client) GeoRadiusByMemberStore(key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoRadiusByMemberStore(ctx, key, member, query)
}

func (cli *Client) GeoSearch(key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoSearch(ctx, key, q)
}

func (cli *Client) GeoSearchLocation(key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoSearchLocation(ctx, key, q)
}

func (cli *Client) GeoSearchStore(key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoSearchStore(ctx, key, store, q)
}

func (cli *Client) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoDist(ctx, key, member1, member2, unit)
}

func (cli *Client) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), cli.options.CtxTimeout)
	defer cancel()
	return cli.Client.GeoHash(ctx, key, members...)
}
