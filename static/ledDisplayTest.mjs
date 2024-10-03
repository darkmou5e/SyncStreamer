import { LedDisplay } from './LedDisplay.mjs'

const led = new LedDisplay(2, 2)
led.mount(document.getElementById("rootEl"))

let i = 0
setInterval(() => {
  const arr = Array(4).fill(null)
  arr[i % 4] = 1
  led.update(arr)
  i++
}, 1000)
