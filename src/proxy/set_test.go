/*
 * Radon
 *
 * Copyright 2018 The Radon Authors.
 * Code is licensed under the GPLv3.
 *
 */

package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelabs/go-mysqlstack/driver"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func TestProxySet(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	fakedbs, proxy, cleanup := MockProxy(log)
	defer cleanup()
	address := proxy.Address()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create table .*", &sqltypes.Result{})
	}

	// create test table.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create table test.t1(id int, b int) partition by hash(id)"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	// set.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		{
			query := "set @@SESSION.radon_streaming_fetch='ON'"
			_, err := client.FetchAll(query, -1)
			assert.Nil(t, err)
		}
		{
			query := "set @@SESSION.radon_streaming_fetch='OFF'"
			_, err := client.FetchAll(query, -1)
			assert.Nil(t, err)
		}
		{
			query := "set @@SESSION.radon_streaming_fetch=true"
			_, err := client.FetchAll(query, -1)
			assert.Nil(t, err)
		}
		{
			query := "set @@SESSION.radon_streaming_fetch=false"
			_, err := client.FetchAll(query, -1)
			assert.Nil(t, err)
		}
		{
			query := "set @@SESSION.radon_streaming_fetch=123"
			_, err := client.FetchAll(query, -1)
			assert.NotNil(t, err)
		}
	}
}
