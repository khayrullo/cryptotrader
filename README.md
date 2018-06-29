# cryptotrader - Coin Exchange Tools

My tools for working with crypto current exchanges.

## Latest Builds

The latest builds from the git master branch can be found at:

https://gitlab.com/crankykernel/cryptotrader/-/jobs/artifacts/master/browse?job=build

## Installation with Go Get

```
go get gitlab.com/crankykernel/cryptotrader
```

## Tools

### GDAX - Ticker

```
cryptotrader gdax ticker [product] ...
```

Make it prettier with `jq`:

```
cryptotrader gdax ticker | jq -c .
```

### KuCoin - Print Trades

```
cryptotrader kucoin --api-key <key> --api-secret <secret> trades
```

### KuCoin - Print Transfers (Deposits and Withdrawals)

```
cryptotrader kucoin transfers
```

Optionally use the KUCOIN_API_KEY and KUCOIN_API_SECRET environment variables
or put them in the configuration file.

## Configuration File Example

Default location: ~/.cryptotrader.yaml

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
