<!-- togo-brand -->
<p align="center"><img src=".github/assets/togo-mark.svg" width="96" alt="togo" /></p>
<h1 align="center">data-bigquery</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a></sub></p>

A togo **data** backend that queries **Google BigQuery** (REST jobs.query).
Registers into [data](https://github.com/togo-framework/data)'s slot.

```bash
togo install togo-framework/data-bigquery
togo provider:use data bigquery
togo config:set BIGQUERY_PROJECT my-gcp-project
# BIGQUERY_TOKEN in .env (secret; OAuth2 access token)
```

MIT © fadymondy
