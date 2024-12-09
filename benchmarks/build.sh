#!/bin/bash

cd ./hotel_hotels
fn init --runtime go hotels
fn build
fn routes create hotel /hotels
fn routes update hotel /hotels -max-concurrency 8 -m 256 --format http --timeout 60s --idle-timeout 600s
cd ..


cd ./hotel_user
fn init --runtime go user
fn build
fn routes create hotel /user
fn routes update hotel /user -max-concurrency 8 -m 256 --format http --timeout 60s --idle-timeout 600s
cd ..

cd ./hotel_recommendations
fn init --runtime go recommendations
fn build
fn routes create hotel /recommendations
fn routes update hotel /recommendations -max-concurrency 8 -m 256 --format http --timeout 60s --idle-timeout 600s
cd ..

cd ./hotel_reservation
fn init --runtime go reservation
fn build
fn routes create hotel /reservation
fn routes update hotel /reservation -max-concurrency 8 -m 256 --format http --timeout 60s --idle-timeout 600s
cd ..