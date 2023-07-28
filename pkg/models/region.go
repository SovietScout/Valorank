package models

type Region struct {
    Region string
    Shard  string
}

func GetRegion(regionStr string) Region {
    // BR & LATAM regions use NA shard
    switch (regionStr) {
    case "br":
        return Region{Region: regionStr, Shard: "na"}
    case "latam":
        return Region{Region: regionStr, Shard: "na"}
    default:
        return Region{Region: regionStr, Shard: regionStr}
    }
}