#!/bin/bash


cd ./original/hotel_user
fn init --runtime go original_user
fn build
fn routes create app /user
fn routes update app /user -max-concurrency 8 -m 256 --timeout 60s
cd ..

cd ./original/hotel_hotel
fn init --runtime go original_hotel
fn build
fn routes create app /hotels
fn routes update app /hotels -max-concurrency 8 -m 256 --timeout 60s
cd ..
