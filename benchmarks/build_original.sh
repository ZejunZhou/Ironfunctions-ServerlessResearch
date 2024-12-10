#!/bin/bash


cd ./original/hotel_user
fn init --runtime go user
fn build
fn routes create app /user
cd ..
