package collections

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wellknown "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"

	"github.com/NibiruChain/nibiru/collections/keys"
)

func deps() (sdk.StoreKey, sdk.Context, codec.BinaryCodec) {
	sk := sdk.NewKVStoreKey("mock")
	dbm := db.NewMemDB()
	ms := store.NewCommitMultiStore(dbm)
	ms.MountStoreWithDB(sk, types.StoreTypeIAVL, dbm)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}
	return sk, sdk.Context{}.WithMultiStore(ms).WithGasMeter(sdk.NewGasMeter(1_000_000_000)), codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
}

func obj(o string) wellknown.BytesValue {
	return wellknown.BytesValue{Value: []byte(o)}
}

func kv(o string) KeyValue[keys.StringKey, wellknown.BytesValue, *wellknown.BytesValue] {
	return KeyValue[keys.StringKey, wellknown.BytesValue, *wellknown.BytesValue]{
		Key:   keys.StringKey(o),
		Value: wellknown.BytesValue{Value: []byte(o)},
	}
}

func TestUpstreamIterAssertions(t *testing.T) {
	// ugly but asserts upstream behavior
	sk, ctx, _ := deps()
	kv := ctx.KVStore(sk)
	kv.Set([]byte("hi"), []byte{})
	i := kv.Iterator(nil, nil)
	err := i.Close()
	require.NoError(t, err)
	require.NoError(t, i.Close())
}

func TestMap(t *testing.T) {
	sk, ctx, cdc := deps()
	m := NewMap[keys.StringKey, wellknown.BytesValue, *wellknown.BytesValue](cdc, sk, 0)

	key := keys.String("id")
	expected := obj("test")

	// test insert and get
	m.Insert(ctx, key, expected)
	got, err := m.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, expected, got)

	// test delete and get error
	err = m.Delete(ctx, key)
	require.NoError(t, err)
	_, err = m.Get(ctx, key)
	require.ErrorIs(t, err, ErrNotFound)

	// test delete errors not exist
	err = m.Delete(ctx, key)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestMap_Iterate(t *testing.T) {
	sk, ctx, cdc := deps()
	m := NewMap[keys.StringKey, wellknown.BytesValue, *wellknown.BytesValue](cdc, sk, 0)

	objs := []KeyValue[keys.StringKey, wellknown.BytesValue, *wellknown.BytesValue]{kv("a"), kv("aa"), kv("b"), kv("bb")}

	m.Insert(ctx, "a", obj("a"))
	m.Insert(ctx, "aa", obj("aa"))
	m.Insert(ctx, "b", obj("b"))
	m.Insert(ctx, "bb", obj("bb"))

	// test iteration ascending
	iter := m.Iterate(ctx, keys.NewRange[keys.StringKey]())
	defer iter.Close()
	for i, o := range iter.All() {
		require.Equal(t, objs[i], o)
	}

	// test iteration descending
	dIter := m.Iterate(ctx, keys.NewRange[keys.StringKey]())
	defer dIter.Close()
	for i, o := range iter.All() {
		require.Equal(t, objs[len(objs)-1-i], o)
	}

	// test all keys

	// test all values
}
