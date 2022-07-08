# URL shortener (backend)
![Build and tests](https://github.com/ibeauregard/url-shortener/actions/workflows/build-and-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibeauregard/url-shortener)](https://goreportcard.com/report/github.com/ibeauregard/url-shortener)
[![codecov](https://codecov.io/gh/ibeauregard/url-shortener/branch/master/graph/badge.svg)](https://codecov.io/gh/ibeauregard/url-shortener)

## Overview

This repo contains the backend components of a URL shortening service similar to what is offered by [TinyURL](https://tinyurl.com/app/) and [Bitly](https://bitly.com/).

The requests for short URLs are made via a simple JSON API. Once you have a short URL, you can use it like any other one in the browser, and you will be redirected to the long URL associated with the shortened one.

## Requirements

Make sure you have [Docker Compose installed](https://docs.docker.com/compose/install/).

## How to run

After you cloned the repo, you can simply execute `make build run` from the project's root directory. When you build the project for the first time, expect a relatively long build, since a lot of dependencies will need to be downloaded. Subsequent builds will be much faster.

See the Makefile for a list of `make` targets that you can use if needed.

Once running, the URL shortener listens and serves on [http://localhost:8080](http://localhost:8080).

## How to use

### Get a short URL
In order to get a short URL, issue a POST to `/api/mappings` with the following body:
```json
{"longUrl": "[any-URL]"}
```
In order to perform the POST request, you can use a tool such as [Postman](https://www.postman.com/). If you wish to use `curl`, run the following command:

```shell
curl --location --request POST 'localhost:8080/api/mappings' \
--header 'Content-Type: application/json' \
--data-raw '{"longUrl": "[any-URL]"}'
```


If both your POST request format and specified long URL are correct, you will receive the following response:

```json
{
    "longUrl": "[long-URL-you-specified]",
    "shortUrl": "[short-url-that-will-redirect-to-long-URL]"
}
```

### Use the short URL

You can use the provided short URL in a web browser, as you would use any other URL. You will be redirected to the long URL that is associated with the short URL.

## How to test

### Unit tests
A unit test suite is automatically run before each build (see Dockerfile). The build will fail is the test suite is not successful.

### Functional tests
You can execute a functional test suite by executing `make func-tests`.

## Key points about design and functionality

### Validation and normalization of URL received

All URLs for which a short URL is requested are validated and normalized. If the long URL does not pass the validation step, no short URL is provided and an error message is sent instead with http status 422 (Unprocessable Entity).

Any received URL is first validated against a custom regular expression. If there is a match, the following steps are subsequently performed:
- add the `http` scheme if none was provided (e.g., `foobar.com` will become `http://foobar.com`)
- remove the port number if that number is the default one associated with the provided scheme (e.g., `http://foobar.com:80` will become `http://foobar.com`).
- convert the host to lowercase
- replace any instances of two or more consecutive slashes with a single slash
- remove any trailing slash
- replace any escape sequences (such as URL encodings) with the escaped characters

### Blacklisting of the application's host

URLs which have the same host and port number as the one on which the application is running are blacklisted (i.e., cannot be associated with a short URL).

For instance, if the application runs on `localhost:8080`, no short URL will be provided for a URL which has `localhost:8080` as its host and port number.

### Short URL key alphabet

"Key" means the sequence of characters that makes up the path in a short URL. For instance, in `locahost:8080/9sg8y`, the key is `9sg8y`.

The alphabet used for generating keys is the following one:
```
23456789BCDFGHJKLMNPQRSTVWXYZbcdfghjkmnpqrstvwxyz-_~!$&=@
```

It consists of the decimal digits and of the uppercase and lowercase letters, plus some special characters. Characters that could cause ambiguity or generate offensive words were removed. This alphabet has a length of 57 characters.

### Key generation

The basis for key generation is the ID of the database row where the associated long and short URL will be stored. When using an autoincrement integer primary key column, SQLite (the chosen DBMS for this project, see [section SQLite](#sqlite) below) provides a handy function to get the highest ID used so far for a given table.

In order to generate the key, the current ID is first converted to its representation using the aforementioned alphabet. This is akin to converting an integer to a decimal representation, except that instead of using "0123456789" as the alphabet, we use the alphabet described in [section Short URL key alphabet](#short-url-key-alphabet) above.

If we stopped there, we would get a functional URL shortener, but it would be pretty easy for anybody to figure out the key alphabet and browse all shortened URLs sequentially.

In order to avoid that, two seemingly random characters are prepended to the key. To accomplish that, we first compute the [CRC-32 checksum](https://en.wikipedia.org/wiki/Cyclic_redundancy_check) of the long URL. We then perform `%alphabetLength^2` on the checksum and convert that modulus to a string representation, which will always be of length 2 (because a padding is performed with the 0th character of the alphabet).

### SQLite

Because the data storage is extremely simple in this project, key value stores were a tempting choice.

A short URL key can easily be converted back to the ID of the record containing that short URL and the associated long URL, providing access to a long URL from a short one for redirect purposes.

More simple, even, the short URL key could have been used as the key and the long URL as the value.

However, since we do not want to store twice the same long URL, we need to be able to search the database with any received long URL to see if it is already stored. So it is obvious that more than one search key is required (both the short URL key and the long URL are used as search keys).

For that reason, key-value stores didn't seem to provide any obvious advantage over SQLite. With SQLite, it is very straightforward to have more than one index key in a table.

### Docker

It was decided to have the application run inside a Docker container. That way, any new developer can get started very quickly without the need to install and/or configure dependencies (e.g., SQLite in this case, but there would likely be more in larger project).
