# cautious-octo-guacamole
## Overview
This endpoint provides up-to-the-minute crypto exchange rates relative to US dollars: https://api.coinbase.com/v2/exchange-rates?currency=USD

That is: each rate is how much of that crypto currency you would get for 1 dollar. So if you received a value for 0.091 for BTC, that's saying it's 0.091 per 1 USD.

Your Task:
You are to make a cli that takes in a USD amount as holdings, and calculates the 70/30 split for 2 given crypto currencies. Stated simply: I have $X I want to keep in BTC and ETH, 70/30 split. How many of each should I buy? An example usage would look like:

input:
example 100 BTC ETH

output:
$70.00 => 0.0025 BTC
$30.00 => 0.0160 ETH
This output tells us: Of our 100$ holdings, 70% of that is 70$, which buys 0.0025 BTC, and 30% of our holdings is 30$, which buys 0.016 ETH.

Out of scope:
To make sure the time is constrained and project doesn't go overboard, feel free to not include tests, we'll check the submissions manually.

Additionally you can always assume a 70/30 split of the amount, and that the arguments will always be in order, e.g. {USD amount} {70% Crypto Ticker} {30% Crypto Ticker}