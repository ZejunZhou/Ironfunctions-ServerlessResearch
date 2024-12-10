#!/bin/bash


cd ./original/hotel_user
fn init --runtime go user
fn build
fn routes create app /user
fn routes update app /user -max-concurrency 8 -m 256 --timeout 60s
cd ..
