import * as fs from 'fs'
import { promisify } from 'util'

const ALPHABET = '0123456789abcdefghijklmnopqrstuvwxyz'

async function countCharOnStream(stream, c) {
    console.time(`read ${c}`)

    return new Promise((resolve) => {
        let count = 0

        stream.on('data', line => {
            line = line
                .toString()
                .toLowerCase()
        
            for (const _c of line) {
                if (c === _c) {
                    count++
                }
            }
        })

        stream.on('end', () => {
            console.timeEnd(`read ${c}`)
            console.log(`${c} ocurrencies: ${count}`)
            resolve(c.repeat(count))
        })
    })
}

for (const c of ALPHABET) {
    countCharOnStream(fs.createReadStream('./in.txt'), c)
}
