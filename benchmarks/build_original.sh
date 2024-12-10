#!/bin/bash


cd ./original/hotel_user
fn init --runtime go original_user
fn build
fn routes create app /user
fn routes update app /user
cd ../..

cd ./original/hotel_hotels
fn init --runtime go original_hotels
fn build
fn routes create app /hotels
fn routes update app /hotels
cd ..
