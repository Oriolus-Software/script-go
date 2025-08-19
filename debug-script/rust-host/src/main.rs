use wasmtime::*;

fn main() {
    let wasm_bytes = std::fs::read("debug-script.wasm").unwrap();
    let engine = Engine::new(
        Config::default()
            .strategy(Strategy::Cranelift)
            .cranelift_opt_level(OptLevel::SpeedAndSize),
    )
    .unwrap();
    let module = Module::new(&engine, &wasm_bytes).unwrap();
    let mut store = Store::new(&engine, ());
    let instance = Instance::new(&mut store, &module, &[]).unwrap();
    let add = instance.get_func(&mut store, "add").unwrap();
    let mut result = [Val::I32(0)];

    for i in 0..10 {
        let start = std::time::Instant::now();
        add.call(&mut store, &[Val::I32(i), Val::I32(i)], &mut result)
            .unwrap();
        let duration = start.elapsed();
        println!("Time taken: {:?}", duration);
    }

    println!("Result: {}", result[0].i32().unwrap());
}
