// TODO: add test
// Internal
// [{i, colorId}]
function calculateDiffOps(oldValue, newValue) {
    const ops = []
    oldValue.forEach((colorId, i) => {
        if (colorId !== newValue[i]) {
            if (newValue[i] === null) {
                ops.push({ i, colorId: null })
            } else {
                ops.push({ i, colorId: newValue[i] })
            }
        }
    })

    return ops
}

const size = 10 // px

// background-color: #240000;
const colors = {
    0: '#FFFFFF', // white
    1: '#2BF0FB', // red
    2: '#00f0ff', // blue
    3: '#65BA00', // green
}

export class LedDisplay {
    _w = 0
    _h = 0

    _currentValue = null // [h * w] colorId
    _mountEl = null
    _pointElements = null  // list [h * w] DOMElements

    constructor(width, height) {
        this._w = width
        this._h = height
    }

    mount(domElement) {
        if (this._currentValue === null) {
            this._init(domElement)
        }
        this._mountEl = domElement
    }

    update(newValue) {
        if (newValue.length !== (this._w * this._h)) {
            throw Error("Size mismatch")
        }
        // [{x, y, colorId}]

        const diffOps = this._calculateDiffOps(newValue)
        this._currentValue = newValue
        this._render(diffOps)
    }

    _widthPx() {
        return this._w * size
    }

    _heightPx() {
        return this._h * size
    }

    _size() {
        return this._w * this._h
    }

    _colOf(i) {
        return i % this._w
    }

    _rowOf(i) {
        return Math.trunc(i / this._w)
    }

    _init(rootEl) {
        this._currentValue = []
        this._pointElements = []

        const svgEl = document.createElementNS("http://www.w3.org/2000/svg", "svg");
        svgEl.setAttribute('viewBox', `0 0 ${this._widthPx()} ${this._heightPx()}`);
        svgEl.setAttribute('width', `100%`);
        svgEl.setAttribute('height', `100%`);
        // svgEl.setAttribute('preserveAspectRatio', `none`);

        for (let i = 0; i < this._size(); i++) {
            const rectEl = document.createElementNS("http://www.w3.org/2000/svg", "rect");
            rectEl.setAttribute('width', size);
            rectEl.setAttribute('height', size);
            rectEl.setAttribute('x', size * this._colOf(i));
            rectEl.setAttribute('y', size * this._rowOf(i));
            rectEl.setAttribute('fill', 'transparent');

            svgEl.appendChild(rectEl)
            this._currentValue.push(null)
            this._pointElements.push(rectEl)
        }

        rootEl.appendChild(svgEl)
    }

    _render(ops) {
        ops.forEach(op => {
            const { i, colorId } = op,
                el = this._pointElements[i]
            if (colorId === null || colorId === 0) {
                el.classList.remove("fade-in")
                el.classList.add("fade-out")
                el.style.opacity = "0"
            } else {
                el.setAttribute('fill', colors[colorId])
                el.classList.remove("fade-out")
                el.classList.add("fade-in")
                el.style.opacity = "1"
            }
            // el.setAttribute('fill', colorId === null ? 'transparent' : colors[colorId])
        })
    }

    // [{x, y, colorId, op: "set"/"clear"}]
    _calculateDiffOps(newValue) {
        return calculateDiffOps(this._currentValue, newValue)
    }
}
