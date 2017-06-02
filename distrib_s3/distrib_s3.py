#!/usr/bin/env python

import argparse
import requests
import time
import boto
from boto.s3.key import Key
import logging
import sys
from logging import getLogger

logger = getLogger(__name__)


def setup_logging():
    root = logging.getLogger()
    root.setLevel(logging.INFO)
    ch = logging.StreamHandler(sys.stdout)
    ch.setLevel(logging.DEBUG)
    formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    ch.setFormatter(formatter)
    root.addHandler(ch)


def strip_slash(s):
    return s.strip('/')


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--opa-url', default='http://localhost:8181/v1', type=strip_slash)
    parser.add_argument('--poll-delay', default=10, type=int)
    parser.add_argument('bucket_name')
    parser.add_argument('filename', type=strip_slash)
    return parser.parse_args()


def main():
    args = parse_args()
    setup_logging()
    logger.info('First line of log stream')
    conn = boto.connect_s3()
    bucket = conn.get_bucket(args.bucket_name)
    while True:
        key = Key(bucket, args.filename)
        content = key.get_contents_as_string()
        t0 = time.time()
        resp = requests.put(args.opa_url + '/policies/' + args.filename, data=content, headers={"Content-Type": "text/plain"})
        if resp.status_code < 200 or resp.status_code >= 300:
            logger.error('Error loading policy into OPA: %s', resp.json())
        else:
            dt = time.time() - t0
            logger.info('Synched policy into OPA (took %.2fms)', dt)
        time.sleep(args.poll_delay)


if __name__ == '__main__':
    main()