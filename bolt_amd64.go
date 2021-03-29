package bolt

// maxMapSize represents the largest mmap size supported by Bolt.
const maxMapSize = 0xFFFFFFFFFFFF // 256TB

// maxAllocSize is the size used when creating array pointers.
const maxAllocSize = 0x7FFFFFFF  // 2^30=1G

// Are unaligned load/stores broken on this arch?
var brokenUnaligned = false
