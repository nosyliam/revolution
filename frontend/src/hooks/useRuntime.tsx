import React, {createContext, useEffect, useState} from "react";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";

interface SetEvent {
    op: "set"
    value: Value
}

interface AppendEvent {
    op: "append"
    primitive: boolean
    keyed: boolean
}

interface DeleteEvent {
    op: "delete"
}

type Event = SetEvent | AppendEvent | DeleteEvent
type Value = number | string | boolean

interface KeyedObject {
    key: string
    object: Object
}

type ListValue = KeyedObject | Object | Value

interface Field {
    value: Value
    dispatch?: React.Dispatch<React.SetStateAction<Value>>
}

interface Reactive {
    // Pass the event data to the object associated with the given path
    Receive(path: Path, event: Event): void

    // Return the reactive object at the given path
    Object(path: Path): Object
}

interface PathComponent {
    val: string
    brackets: boolean
}

export class Path extends String {
    private index: number = 0;
    components: PathComponent[] = [];

    constructor(path: string, components?: PathComponent[]) {
        super(path);
        if (components) {
            this.components = components
            return
        }
        const regex = /(\w+)|\[(.*?)\]/g;

        let match;
        while ((match = regex.exec(path)) !== null) {
            if (match[1]) {
                this.components = [...this.components, {val: match[1], brackets: false}]
            }
            if (match[2]) {
                this.components = [...this.components, {val: match[2], brackets: true}]
            }
        }
    }

    public reset(): Path {
        this.index = 0
        return this
    }

    public increment(): Path {
        this.index += 1
        return this
    }

    public decrement(): Path {
        this.index -= 1
        return this
    }


    public get value(): string {
        return this.components[this.index].val
    }

    public get final(): boolean {
        return this.components.length == this.index
    }

    public get nextFinal(): boolean {
        return this.components.length - 1 == this.index
    }


    public extend(path: string, brackets?: boolean): Path {
        return new Path(brackets ? `${this}[${path}]` : `${this}.${path}`, [...this.components, {
            val: path,
            brackets: Boolean(brackets)
        }])
    }
}

export class Object implements Reactive {
    private readonly path: Path;
    private readonly runtime: Runtime;

    private objects: {[field: string]: Reactive} = {};
    private values: {[field: string]: Field} = {};

    constructor(path: Path, runtime: Runtime) {
        this.runtime = runtime
        this.path = path
    }

    public Value<T extends Value>(field: string, defaultValue: T): T {
        let def: T = defaultValue
        if (this.values[field]) {
            def = this.values[field].value as T
            console.log("found value", String(this.path), def)
        } else {
            console.log("setting value", String(this.path), def)
        }
        const [value, dispatch] = useState<Value>(def)
        this.values[field] = {value: def, dispatch: dispatch}
        return value as T
    }

    public Set<T extends Value>(field: string, value: T) {
        let data: Field = this.values[field] || {value: value}
        if (data.dispatch)
            data.dispatch(value)
        this.values[field] = {...data, value: value}
        this.runtime.Emit(this.Field(field), {op: "set", value: value})
    }

    public List<T extends ListValue>(field: string): List<T> {
        if (!this.objects[field]) {
            this.objects[field] = new List(this.path.extend(field), this.runtime)
        }
        return this.objects[field] as List<T>
    }

    public Receive(path: Path, event: Event): void {
        switch (event.op) {
            case "set":
                if (!path.nextFinal) {
                    if (!this.objects[path.value]) {
                        console.log("creating object", String(path), path.value)
                        this.objects[path.value] = new Object(this.path.extend(path.value), this.runtime)
                    }
                    this.objects[path.value].Receive(path.increment(), event)
                } else {
                    let field: Field | undefined
                    if ((field = this.values[path.value]) != undefined) {
                        field.value = event.value
                        if (field.dispatch) {
                            console.log("dispatching received", String(path), event.value)
                            field.dispatch(event.value)
                        }
                    } else {
                        this.values[path.value] = {value: event.value}
                    }
                }
                break
            case "append":
                path.increment()
                console.log("appending", String(path), String(this.path), path.final, path.nextFinal)
                if (path.nextFinal) {
                    path.decrement()
                    if (!this.objects[path.value])
                        this.objects[path.value] = new List(this.path.extend(path.value), this.runtime)
                    let list = this.objects[path.value] as List<any>
                    list.primitive = event.primitive
                    list.keyed = event.keyed
                    console.log("set list", String(path), String(this.path), path.value)
                    path.increment()
                }
                this.objects[path.decrement().value].Receive(path.increment(), event)
                break
            case "delete":
                this.objects[path.increment().value].Receive(path, event)
        }
    }

    Object(path: Path): Object {
        if (path.final) {
            return this
        }
        if (!this.objects[path.value]) {
            this.objects[path.value] = new Object(this.path.extend(path.value), this.runtime)
        }
        return this.objects[path.value].Object(path.increment())
    }

    public Field(field: string): Path {
        return new Path(`${this.path}.${field}`)
    }
}

export class List<T extends ListValue> implements Reactive {
    private readonly path: Path;
    private readonly runtime: Runtime;

    private dispatch?: React.Dispatch<React.SetStateAction<T[]>>
    private values: T[] = []

    public primitive: boolean | undefined = undefined;
    public keyed: boolean = false;

    constructor(path: Path, runtime: Runtime) {
        this.path = path;
        this.runtime = runtime;
    }

    private index(path: Path): number {
        let index: number;
        if (this.keyed) {
            index = (this.values as KeyedObject[]).findIndex((v) => v.key == path.value)
        } else {
            index = Number(path.value)
        }
        return index
    }

    public Values(): T[] {
        const [values, dispatch] = useState<T[]>([])
        this.dispatch = dispatch
        useEffect(() => {
            return () => {
                this.dispatch = undefined
            }
        })
        return values as T[]
    }

    public Delete(key: string | number) {
        this.runtime.Emit(this.Key(key), {op: "delete"})
    }

    Receive(path: Path, event: Event): void {
        const index = this.index(path)
        switch (event.op) {
            case "set":
                if (this.primitive) {
                    (this.values[index] as Value) = event.value
                } else if (this.keyed) {
                    console.log("sending2", String(path), event, this.values, this.keyed);
                    (this.values[index] as KeyedObject).object.Receive(path.increment(), event)
                } else {
                    console.log("sending", String(path), event, this.values, this.keyed);
                    (this.values[index] as Object).Receive(path.increment(), event)
                }
                break;
            case "append":
                if (this.keyed) {
                    (this.values as KeyedObject[]) = [...(this.values as KeyedObject[]), {
                        key: path.value,
                        object: new Object(path, this.runtime)
                    }]
                    console.log("updated values", path, [...this.values])
                } else if (!this.primitive) {
                    (this.values as Object[]) = [...(this.values as Object[]), new Object(path, this.runtime)]
                }
                console.log("appended", this.keyed, this.values)
                if (this.dispatch)
                    this.dispatch(this.values)
            case "delete":
                this.values = [...this.values.splice(index)]
                if (this.dispatch)
                    this.dispatch([...this.values.splice(index)])
        }
    }

    Object(path: Path): Object {
        const index = this.index(path)
        return (this.values[index] as Object).Object(path.increment())
    }

    public Key(field: string | number): Path {
        return new Path(`${this.path}[${field}]`)
    }
}

export class Runtime {
    static RootNames: string[] = ["settings", "state", "database"]
    private roots: {[name: string]: Reactive} = {};

    constructor() {
        for (const root of Runtime.RootNames) {
            this.roots[root] = new Object(new Path(root), this)
        }
        EventsOn("set", (path: string, value: Value) => this.Receive(new Path(path), {
            op: "set",
            value: value
        }))
        EventsOn("append", (path: string, primitive: boolean, keyed: boolean) => this.Receive(new Path(path), {
            op: "append",
            primitive: primitive,
            keyed: keyed
        }))
        EventsOn("delete", (path: string, value: Value) => this.Receive(new Path(path), {op: "delete"}))
        // @ts-ignore
        window.dataRuntime = this
        // @ts-ignore
        window.pathObject = Path
    }

    Receive(path: Path, event: Event): void {
        this.roots[path.value].Receive(path.increment(), event)
    }

    public Emit(path: Path, event: Event) {
        EventsEmit(`${event.op}_client`, String(path), ...(event.op == 'set' ? [event.value] : []))
    }

    public Ready() {
        EventsEmit("ready")
    }

    public Object(path: Path | string): Object {
        if (typeof path == 'string') {
            path = new Path(path as string)
        }
        path.reset()
        return this.roots[path.value].Object(path.increment())
    }
}

export const RuntimeContext = createContext<Runtime>(new Runtime())