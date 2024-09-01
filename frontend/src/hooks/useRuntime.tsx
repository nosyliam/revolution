import React, {createContext, useEffect, useState} from "react";
import {EventsEmit} from "../../wailsjs/runtime";

interface SetEvent {
    op: "set"
    value: Value
}

interface AppendEvent {
    op: "append"
    primitive?: boolean
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
                this.components = [...this.components, {val: match[1], brackets: true}]
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

    public get value(): string {
        return this.components[this.index].val
    }

    public get final(): boolean {
        return this.components.length - 1 == this.index
    }

    public root(): string {
        return this.components[0].val
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

    private objects: Map<string, Reactive>;
    private values: Map<string, Field>;

    constructor(path: Path, runtime: Runtime) {
        this.runtime = runtime
        this.path = path
        this.values = new Map();
        this.objects = new Map();
    }

    public Value<T extends Value>(field: string, defaultValue: T): T {
        let def: T = defaultValue
        if (this.values.has(field)) {
            def = this.values.get(field)!.value as T
        }
        const [value, dispatch] = useState<Value>(def)
        this.values.set(field, {value: def, dispatch: dispatch})
        return value as T
    }

    public Set<T extends Value>(field: string, value: T) {
        let data: Field = this.values.get(field) || {value: value}
        if (data.dispatch)
            data.dispatch(value)
        this.values.set(field, {...data, value: value})
    }

    public List(field: string): List {
        if (!this.objects.has(field)) {
            this.objects.set(field, new List(this.path.extend(field), this.runtime))
        }
        return this.objects.get(field)! as List
    }

    public Receive(path: Path, event: Event): void {
        let object: Reactive | undefined
        switch (event.op) {
            case "set":
                if (!path.final) {
                    if (!this.objects.has(path.value)) {
                        this.objects.set(path.value, new Object(this.path.extend(path.value), this.runtime))
                    }
                    this.objects.get(path.value)!.Receive(path.increment(), event)
                } else {
                    let field: Field | undefined
                    if ((field = this.values.get(path.value)) != undefined) {
                        field.value = event.value
                        if (field.dispatch)
                            field.dispatch(event.value)
                    } else {
                        this.values.set(path.value, {value: event.value})
                    }
                }
                break
            case "append":
                if (!this.objects.has(path.value))
                    this.objects.set(path.value, new List(this.path.extend(path.value), this.runtime))
                let list = this.objects.get(path.value)! as List
                list.primitive = event.primitive
            case "delete":
                this.objects.get(path.value)!.Receive(path.increment(), event)
        }
    }

    Object(path: Path): Object {
        if (path.final)
            return this
        if (!this.objects.has(path.value))
            this.objects.set(path.value, new Object(this.path.extend(path.value), this.runtime))
        return this.objects.get(path.value)!.Object(path.increment())
    }

    public Field(field: string): Path {
        return new Path(`${this.path}.${field}`)
    }
}

export class List implements Reactive {
    private readonly path: Path;
    private readonly runtime: Runtime;

    private dispatch?: React.Dispatch<React.SetStateAction<ListValue[]>>
    private values: ListValue[] = []

    public primitive: boolean | undefined = undefined;
    public keyed: boolean = false;

    constructor(path: Path, runtime: Runtime) {
        this.path = path;
        this.runtime = runtime;
    }

    public Values<T extends ListValue>(): T[] {
        const [values, dispatch] = useState<ListValue[]>([])
        this.dispatch = dispatch
        useEffect(() => {
            return () => {
                this.dispatch = undefined
            }
        })
        return values as T[]
    }

    Receive(path: Path, event: Event): void {
        if (!this.values)
            [this.values, this.dispatch] = useState<ListValue[]>([])
        let index: number;
        if (this.keyed) {
            index = (this.values as KeyedObject[]).findIndex((v) => v.key == path.value)
        } else {
            index = Number(path.value)
        }
        switch (event.op) {
            case "set":
                if (this.primitive) {
                    this.values[index] = event.value
                } else {
                    (this.values[index] as Object).Receive(path.extend(String(index), true), event)
                }
                break;
            case "append":
                this.values = [...this.values, new Object(path, this.runtime)]
                if (this.dispatch)
                    this.dispatch(this.values)
            case "delete":
                this.values = [...this.values.splice(index)]
                if (this.dispatch)
                    this.dispatch(this.values)
        }
    }

    Object(path: Path): Object {
        throw new Error("Method not implemented.");
    }

    public Key(field: string | number): Path {
        return new Path(`${this.path}[${field}]`)
    }
}

export class Runtime {
    static rootNames: string[] = ["settings", "state", "database"]
    private roots: Map<string, Reactive>;

    constructor() {
        this.roots = new Map();
    }

    public Emit(path: Path, event: Event) {
    }

    public Object(path: Path): Object {
    }
}

export const RuntimeContext = createContext<Runtime>(new Runtime())