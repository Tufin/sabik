# sabik

[![CircleCI](https://circleci.com/gh/Tufin/sabik.svg?style=svg&circle-token=7d3fae746065ab6be0acd5b2ef9a2fd9c2bfd111)](https://circleci.com/gh/Tufin/sabik)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/sabik)](https://goreportcard.com/report/github.com/tufin/sabik)
<!-- [![codecov](https://codecov.io/gh/Tufin/sabik/branch/master/graph/badge.svg)](https://codecov.io/gh/Tufin/sabik) -->

Use the follow environment variables to setup sabik in a service name _customer_ on _generic-bank/retail_ for example:
```
- name: TUFIN_DOMAIN
  value: generic-bank
- name: TUFIN_PROJECT
  value: retail
- name: TUFIN_SABIK_ENABLE
  value: "true"
- name: TUFIN_SABIK_SERVICE_NAME
  value: customer
- name: TUFIN_SABIK_URL
  value: https://validator-hardcoded-spec-xiixymmvca-ew.a.run.app
```
