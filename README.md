# ctrader - Coin Exchange Tools

My tools for working with crypto current exchanges.

## Installation with Go Get

```
go get gitlab.com/crankykernel/ctrader
```

## Tools

### GDAX - Ticker

```
ctrader gdax ticker [product] ...
```

Make it prettier with `jq`:

```
ctrader gdax ticker | jq -c .
```

### KuCoin - List Trades

```
ctrader kucoin --api-key <key> --api-secret <secret> trades
```

Optionally use the KUCOIN_API_KEY and KUCOIN_API_SECRET environment variables
or put them in the configuration file.

## Configuration File Example

Default location: ~/.ctrader.yaml

```
# KuCoin API key
kucoin.api.key: xxx

# KuCoin API secret
kucoin.api.secret: xxx

# QuadrigaCX
quadriga.api.client-id: xxx
quadriga.api.key: xxx
quadriga.api.secret: xxx

# Kraken
kraken.api.key: xxx
kraken.api.secret: xxx
```