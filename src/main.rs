use std::time::Instant;

use wasmtime::*;

fn main() {
    let engine = Engine::default();
    let module = Module::new(&engine, include_bytes!("../script/main.wasm")).unwrap();
    let mut store = Store::new(&engine, 1024);

    let mut linker = Linker::new(&engine);
    let now_func = Func::wrap(&mut store, || -> i64 {
        Instant::now().elapsed().as_nanos() as i64
    });
    linker.define(&mut store, "time", "now", now_func).unwrap();

    let instance = linker.instantiate(&mut store, &module).unwrap();
    let func = instance.get_func(&mut store, "add").unwrap();
    let mut results = [Val::I32(0)];

    for _ in 0..100 {
        let now = Instant::now();
        func.call(&mut store, &[Val::I32(1), Val::I32(2)], &mut results)
            .unwrap();
        let elapsed = now.elapsed().as_nanos();
        println!("Time: {}", elapsed);
    }
}
