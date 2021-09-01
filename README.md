# ziprox

Module `ziprox` is a library for finding nearby zip codes.

## Overview

The `ziprox.Map` structure groups zip codes into buckets of increasingly large
distances, for every zip code. Those buckets are hard-coded as distances of 5,
10, 20, 50, 100, 200, and 500 miles. The underlying dataset comes from free `.csv`
files provided by NBER.

## NBER

All zip code distance data is from National Bureau of Economic Research. [NBER](https://www.nber.org)

Distance data is available in multiple [datasets](https://www.nber.org/research/data/zip-code-distance-database).

`ziprox` works with `.csv` files for

 - [25 miles](https://nber.org/distance/2016/gaz/zcta5/gaz2016zcta5distance25miles.csv.zip)
 -- uses about 42 MB in memory
 - [50 miles](https://nber.org/distance/2016/gaz/zcta5/gaz2016zcta5distance50miles.csv.zip)
 -- uses about 86 MB in memory
 - [100 miles](https://nber.org/distance/2016/gaz/zcta5/gaz2016zcta5distance100miles.csv.zip)
 -- uses about 144 MB in memory
 - [500 miles](https://nber.org/distance/2016/gaz/zcta5/gaz2016zcta5distance500miles.csv.zip)
 -- uses about 1.5 GB in memory

## Usage

```golang
import "gophers.dev/pkgs/ziprox"

# ...

f, _ := os.Open("zips25.csv")
db, _ := zips.New(f)
origin := ziprox.Zip(75234)
nearby := db.Within(origin, 5)

# nearby = [75229 75006 75244 75220 75039 75001]
```
