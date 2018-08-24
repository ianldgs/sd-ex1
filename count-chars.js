const { createReadStream } = require('fs')

const char = process.argv[2]

console.log('read', char)
console.time(`read ${char}`)

const stream = createReadStream('./in.txt')

let count = 0

stream.on('data', line => {
    line = line
        .toString()
        .toLowerCase()

    for (const c of line) {
        if (char === c) {
            count++
        }
    }
})

stream.on('end', () => {
    console.timeEnd(`read ${char}`)

    process.send(char.repeat(count))
})
