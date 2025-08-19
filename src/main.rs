use std::time::Instant;

use wasmtime::*;

fn main() {
    let engine = Engine::new(Config::new().strategy(Strategy::Cranelift)).unwrap();
    let module = Module::new(&engine, std::fs::read("./script/main.wasm").unwrap()).unwrap();
    let mut store = Store::new(&engine, 1024);

    let mut linker = Linker::new(&engine);
    let now_func = Func::wrap(&mut store, || -> i64 {
        Instant::now().elapsed().as_nanos() as i64
    });
    linker.define(&mut store, "time", "now", now_func).unwrap();

    let instance = linker.instantiate(&mut store, &module).unwrap();
    // let func = instance.get_func(&mut store, "add").unwrap();
    let allocate_func = instance.get_func(&mut store, "allocate").unwrap();
    let deallocate_func = instance.get_func(&mut store, "deallocate").unwrap();
    let mut results = [Val::I32(0)];

    // for _ in 0..1000 {
    //     let start = Instant::now();
    //     allocate_func
    //         .call(&mut store, &[Val::I32(128)], &mut results)
    //         .unwrap();
    //     let ptr = results[0].i32().unwrap();
    //     let memory = instance.get_memory(&mut store, "memory").unwrap();

    //     deallocate_func
    //         .call(&mut store, &[Val::I32(ptr)], &mut [])
    //         .unwrap();

    //     let elapsed = start.elapsed().as_nanos();

    //     println!(
    //         "Data: {} {}",
    //         format_bytes(memory.data_size(&mut store)),
    //         elapsed,
    //     );
    // }

    // for _ in 0..100 {
    //     let now = Instant::now();
    //     func.call(&mut store, &[Val::I32(1), Val::I32(2)], &mut results)
    //         .unwrap();
    //     let elapsed = now.elapsed().as_nanos();
    //     println!("Time: {}", elapsed);
    // }
}

fn format_bytes(bytes: usize) -> String {
    if bytes < 1024 {
        return format!("{} B", bytes);
    } else if bytes < 1024 * 1024 {
        return format!("{} KB", bytes / 1024);
    } else if bytes < 1024 * 1024 * 1024 {
        return format!("{} MB", bytes / (1024 * 1024));
    } else if bytes < 1024 * 1024 * 1024 * 1024 {
        return format!("{} GB", bytes / (1024 * 1024 * 1024));
    }
    format!("{} TB", bytes / (1024 * 1024 * 1024 * 1024))
}
