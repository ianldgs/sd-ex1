import * as fs from 'fs'
import { promisify } from 'util'

import { fork } from 'child_process'

const ALPHABET = '0123456789abcdefghijklmnopqrstuvwxyz'

for (const c of ALPHABET) {
    const child = fork('count-chars.js', [c])
    child.on('message', text => {
        console.log(c, text.length)
    })
}
