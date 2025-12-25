extends "res://test_base.gd"

var custom_signal_emitted = null

class TestClass:
	func test(p_msg: String) -> String:
		return p_msg + " world"

func _ready():
	var example: Example = $Example
	var physics: PhysicsValidation = $PhysicsValidation
	var input_probe: InputProbe = $InputProbe
	if OS.has_environment("GODOT_GO_LEAK_TEST"):
		await run_leak_test(example)
		exit_with_status()
		return
	test_suite(1, example)
	await physics_test_suite(physics)
	await input_test_suite(input_probe)
	# example.group_subgroup_custom_position = Vector2(0, 0)
	# custom_signal_emitted = null
	# var t = get_tree()
	# if t != null:
	# 	await t.create_timer(3.0).timeout
	# test_suite(2, example)
	exit_with_status()

func run_leak_test(example: Example) -> void:
	var duration_seconds := get_env_int("GODOT_GO_LEAK_TEST_SECONDS", 600)
	var interval_ms := get_env_int("GODOT_GO_LEAK_TEST_INTERVAL_MS", 100)
	var iterations := get_env_int("GODOT_GO_LEAK_TEST_ITERATIONS", 1000)
	var max_heap_bytes := get_env_int("GODOT_GO_LEAK_TEST_MAX_HEAP_BYTES", 10 * 1024 * 1024)
	var max_heap_objects := get_env_int("GODOT_GO_LEAK_TEST_MAX_HEAP_OBJECTS", 5000)

	print("leak test: duration=%ss interval_ms=%s iterations=%s max_heap_bytes=%s max_heap_objects=%s" % [
		duration_seconds, interval_ms, iterations, max_heap_bytes, max_heap_objects
	])

	example.leak_check_start()

	var end_time := Time.get_ticks_msec() + duration_seconds * 1000
	while Time.get_ticks_msec() < end_time:
		example.leak_check_tick(iterations)
		await get_tree().create_timer(float(interval_ms) / 1000.0).timeout

	var results := example.leak_check_finish(max_heap_bytes, max_heap_objects)
	var ok := results.get("ok", false)
	print("leak test results: ", results)
	assert_true(ok)

func get_env_int(name: String, default_value: int) -> int:
	if OS.has_environment(name):
		var raw := OS.get_environment(name)
		if raw != "":
			return int(raw)
	return default_value

func test_suite(i: int, example: Example):
	print("test suite run %d" % [i])
	# Signal.
	example.emit_custom_signal("Button", 42)
	assert_equal(custom_signal_emitted, ["Button", 42])

	# To string.
	assert_equal(example.to_string(),'[ GDExtension::Example <--> Instance ID:%s ]' % example.get_instance_id())
	# It appears there's a bug with instance ids :-(
	#assert_equal($Example/ExampleMin.to_string(), 'ExampleMin:[Wrapped:%s]' % $Example/ExampleMin.get_instance_id())

	# godot-go will probably not support static methods since they don't exist in go
	# Call static methods.
	# assert_equal(Example.test_static(9, 100), 109);
	# It's void and static, so all we know is that it didn't crash.
	# Example.test_static2()

	# Property list.
	example.property_from_list = Vector3(100, 200, 300)
	assert_equal(example.property_from_list, Vector3(100, 200, 300))
	var prop_list = example.get_property_list()
	for prop_info in prop_list:
		if prop_info['name'] == 'mouse_filter':
			assert_equal(prop_info['usage'], PROPERTY_USAGE_NO_EDITOR)

	# Call simple methods.
	example.simple_func()
	assert_equal(custom_signal_emitted, ['simple_func', 3])
	example.simple_const_func(123)
	assert_equal(custom_signal_emitted, ['simple_const_func', 4])

	# Pass custom reference.
	# assert_equal(example.custom_ref_func(null), -1)
	# var ref1 = ExampleRef.new()
	# ref1.id = 27
	# assert_equal(example.custom_ref_func(ref1), 27)
	# ref1.id += 1;
	# assert_equal(example.custom_const_ref_func(ref1), 28)

	# Pass core reference.
	assert_equal(example.image_ref_func(null), "invalid")
	# assert_equal(example.image_const_ref_func(null), "invalid")
	var image = Image.new()
	assert_equal(example.image_ref_func(image), "valid")
	# assert_equal(example.image_const_ref_func(image), "valid")

	# Return values.
	assert_equal(example.return_something("some string", 7.0/6, 7.0/6 * 1000, 2147483647, -127, -32768, 2147483647, 9223372036854775807), "1. some string42, 2. %.6f, 3. %f, 4. 2147483647, 5. -127, 6. -32768, 7. 2147483647, 8. 9223372036854775807" % [7.0/6, 7.0/6 * 1000])
	assert_equal(example.return_something_const(), get_viewport())
	var null_ref = example.return_empty_ref()
	assert_equal(null_ref, null)
	# var ret_ref = example.return_extended_ref()
	# assert_not_equal(ret_ref.get_instance_id(), 0)
	# assert_equal(ret_ref.get_id(), 0)
	assert_equal(example.get_v4(), Vector4(1.2, 3.4, 5.6, 7.8))
	assert_equal(example.test_node_argument(example), example)

	# VarArg method calls.
	# var var_ref = ExampleRef.new()
	# assert_not_equal(example.extended_ref_checks(var_ref).get_instance_id(), var_ref.get_instance_id())
	assert_equal(example.varargs_func("some", "arguments", "to", "test"), 4)
	assert_equal(example.varargs_func("some"), 1)
	assert_equal(example.varargs_func_nv("some", "arguments", "to", "test"), 46)
	example.varargs_func_void("some", "arguments", "to", "test")
	assert_equal(custom_signal_emitted, ["varargs_func_void", 5])

	# Method calls with default values.
	assert_equal(example.def_args(), 300)
	assert_equal(example.def_args(50), 250)
	assert_equal(example.def_args(50, 100), 150)

	# Array and Dictionary
	assert_equal(example.test_array(), [1, 2])
	assert_equal(example.test_tarray(), [ Vector2(1, 2), Vector2(2, 3) ])
	assert_equal(example.test_dictionary(), {"hello": "world", "foo": "bar"})
	var array: Array[int] = [1, 2, 3]
	print("array: ", array)
	assert_equal(example.test_tarray_arg(array), 6)

	example.callable_bind()
	assert_equal(custom_signal_emitted, ["bound", 11])

	# String += operator
	assert_equal(example.test_string_ops(), "ABCĎE")

	# UtilityFunctions::str()
	assert_equal(example.test_str_utility(), "Hello, World! The answer is 42")

	# UtilityFunctions::instance_from_id()
	assert_equal(example.test_instance_from_id_utility(), example)

	# # Test converting string to char* and doing comparison.
	# assert_equal(example.test_string_is_fourty_two("blah"), false)
	# assert_equal(example.test_string_is_fourty_two("fourty two"), true)

	# # String::resize().
	# assert_equal(example.test_string_resize("What"), "What!?")

	# mp_callable() with void method.
	# var mp_callable: Callable = example.test_callable_mp()
	# assert_equal(mp_callable.is_valid(), true)
	# mp_callable.call(example, "void", 36)
	# assert_equal(custom_signal_emitted, ["unbound_method1: Example - void", 36])

	# # Check that it works with is_connected().
	# assert_equal(example.renamed.is_connected(mp_callable), false)
	# example.renamed.connect(mp_callable)
	# assert_equal(example.renamed.is_connected(mp_callable), true)
	# # Make sure a new object is still treated as equivalent.
	# assert_equal(example.renamed.is_connected(example.test_callable_mp()), true)
	# assert_equal(mp_callable.hash(), example.test_callable_mp().hash())
	# example.renamed.disconnect(mp_callable)
	# assert_equal(example.renamed.is_connected(mp_callable), false)

	# # mp_callable() with return value.
	# var mp_callable_ret: Callable = example.test_callable_mp_ret()
	# assert_equal(mp_callable_ret.call(example, "test", 77), "unbound_method2: Example - test - 77")

	# # mp_callable() with const method and return value.
	# var mp_callable_retc: Callable = example.test_callable_mp_retc()
	# assert_equal(mp_callable_retc.call(example, "const", 101), "unbound_method3: Example - const - 101")

	# # mp_callable_static() with void method.
	# var mp_callable_static: Callable = example.test_callable_mp_static()
	# mp_callable_static.call(example, "static", 83)
	# assert_equal(custom_signal_emitted, ["unbound_static_method1: Example - static", 83])

	# # Check that it works with is_connected().
	# assert_equal(example.renamed.is_connected(mp_callable_static), false)
	# example.renamed.connect(mp_callable_static)
	# assert_equal(example.renamed.is_connected(mp_callable_static), true)
	# # Make sure a new object is still treated as equivalent.
	# assert_equal(example.renamed.is_connected(example.test_callable_mp_static()), true)
	# assert_equal(mp_callable_static.hash(), example.test_callable_mp_static().hash())
	# example.renamed.disconnect(mp_callable_static)
	# assert_equal(example.renamed.is_connected(mp_callable_static), false)

	# # mp_callable_static() with return value.
	# var mp_callable_static_ret: Callable = example.test_callable_mp_static_ret()
	# assert_equal(mp_callable_static_ret.call(example, "static-ret", 84), "unbound_static_method2: Example - static-ret - 84")

	# # CallableCustom.
	# var custom_callable: Callable = example.test_custom_callable();
	# assert_equal(custom_callable.is_custom(), true);
	# assert_equal(custom_callable.is_valid(), true);
	# assert_equal(custom_callable.call(), "Hi")
	# assert_equal(custom_callable.hash(), 27);

func physics_test_suite(physics: PhysicsValidation) -> void:
	print("physics test suite run")
	var rig := setup_physics_rig()
	await get_tree().physics_frame

	assert_true(physics.enable_ccd(rig.ball))

	var material := PhysicsMaterial.new()
	assert_true(physics.configure_material(rig.ball, material, 0.2, 0.8))

	physics.reset_area_counts()
	assert_true(physics.bind_area(rig.trigger))

	physics.apply_flipper_impulse(rig.ball, Vector2(0, 400), Vector2.ZERO)
	await wait_physics_frames(10)

	assert_true(physics.get_linear_speed(rig.ball) > 0.1)
	assert_true(physics.get_area_enter_count() > 0)

	var node_a := str(rig.joint.get_path_to(rig.anchor))
	var node_b := str(rig.joint.get_path_to(rig.flipper))
	assert_true(physics.configure_pin_joint(rig.joint, node_a, node_b, -0.5, 0.5, 0.2, 8.0))

func setup_physics_rig() -> Dictionary:
	var rig_root := Node2D.new()
	rig_root.name = "PhysicsRig"
	add_child(rig_root)

	var ball := RigidBody2D.new()
	ball.name = "Ball"
	ball.position = Vector2(200, 60)
	ball.contact_monitor = true
	ball.max_contacts_reported = 4
	var ball_shape := CollisionShape2D.new()
	var ball_circle := CircleShape2D.new()
	ball_circle.radius = 6.0
	ball_shape.shape = ball_circle
	ball.add_child(ball_shape)
	rig_root.add_child(ball)

	var floor := StaticBody2D.new()
	floor.name = "Floor"
	floor.position = Vector2(200, 260)
	var floor_shape := CollisionShape2D.new()
	var floor_rect := RectangleShape2D.new()
	floor_rect.size = Vector2(240, 12)
	floor_shape.shape = floor_rect
	floor.add_child(floor_shape)
	rig_root.add_child(floor)

	var trigger := Area2D.new()
	trigger.name = "Trigger"
	trigger.position = Vector2(200, 180)
	trigger.monitoring = true
	trigger.monitorable = true
	var trigger_shape := CollisionShape2D.new()
	var trigger_rect := RectangleShape2D.new()
	trigger_rect.size = Vector2(80, 12)
	trigger_shape.shape = trigger_rect
	trigger.add_child(trigger_shape)
	rig_root.add_child(trigger)

	var anchor := StaticBody2D.new()
	anchor.name = "FlipperAnchor"
	anchor.position = Vector2(60, 220)
	var anchor_shape := CollisionShape2D.new()
	var anchor_rect := RectangleShape2D.new()
	anchor_rect.size = Vector2(10, 10)
	anchor_shape.shape = anchor_rect
	anchor.add_child(anchor_shape)
	rig_root.add_child(anchor)

	var flipper := RigidBody2D.new()
	flipper.name = "Flipper"
	flipper.position = Vector2(80, 220)
	flipper.gravity_scale = 0.0
	var flipper_shape := CollisionShape2D.new()
	var flipper_rect := RectangleShape2D.new()
	flipper_rect.size = Vector2(40, 8)
	flipper_shape.shape = flipper_rect
	flipper.add_child(flipper_shape)
	rig_root.add_child(flipper)

	var joint := PinJoint2D.new()
	joint.name = "FlipperJoint"
	joint.position = anchor.position
	rig_root.add_child(joint)

	return {
		"root": rig_root,
		"ball": ball,
		"floor": floor,
		"trigger": trigger,
		"anchor": anchor,
		"flipper": flipper,
		"joint": joint,
	}

func wait_physics_frames(count: int) -> void:
	for _i in range(count):
		await get_tree().physics_frame

func input_test_suite(probe: InputProbe) -> void:
	print("input test suite run")
	probe.reset_counts()
	probe.set_handle_input(false)
	await get_tree().process_frame
	send_key_events(KEY_SPACE, 20)
	await get_tree().process_frame
	assert_equal(probe.get_input_count(), 20)
	assert_equal(probe.get_unhandled_count(), 20)

	probe.reset_counts()
	probe.set_handle_input(true)
	await get_tree().process_frame
	send_key_events(KEY_SPACE, 15)
	await get_tree().process_frame
	assert_equal(probe.get_input_count(), 15)
	assert_equal(probe.get_unhandled_count(), 0)

func send_key_events(keycode: Key, count: int) -> void:
	for _i in range(count):
		var ev := InputEventKey.new()
		ev.keycode = keycode
		ev.pressed = true
		Input.parse_input_event(ev)
	# assert_equal(custom_callable.get_object(), null);
	# assert_equal(custom_callable.get_method(), "");
	# assert_equal(str(custom_callable), "<MyCallableCustom>");

	# PackedArray iterators
	assert_equal(example.test_vector_ops(), 105)
	# assert_equal(example.test_vector_init_list(), 105)

	# Properties.
	assert_equal(example.group_subgroup_custom_position, Vector2(0, 0))
	example.group_subgroup_custom_position = Vector2(50, 50)
	assert_equal(example.group_subgroup_custom_position, Vector2(50, 50))

	# # Test Object::cast_to<>() and that correct wrappers are being used.
	# var control = Control.new()
	# var sprite = Sprite2D.new()
	# var example_ref = ExampleRef.new()

	# assert_equal(example.test_object_cast_to_node(control), true)
	# assert_equal(example.test_object_cast_to_control(control), true)
	# assert_equal(example.test_object_cast_to_example(control), false)

	# assert_equal(example.test_object_cast_to_node(example), true)
	# assert_equal(example.test_object_cast_to_control(example), true)
	# assert_equal(example.test_object_cast_to_example(example), true)

	# assert_equal(example.test_object_cast_to_node(sprite), true)
	# assert_equal(example.test_object_cast_to_control(sprite), false)
	# assert_equal(example.test_object_cast_to_example(sprite), false)

	# assert_equal(example.test_object_cast_to_node(example_ref), false)
	# assert_equal(example.test_object_cast_to_control(example_ref), false)
	# assert_equal(example.test_object_cast_to_example(example_ref), false)

	# control.queue_free()
	# sprite.queue_free()

	# Test conversions to and from Variant.
	# assert_equal(example.test_variant_vector2i_conversion(Vector2i(1, 1)), Vector2i(1, 1))
	# assert_equal(example.test_variant_vector2i_conversion(Vector2(1.0, 1.0)), Vector2i(1, 1))
	# assert_equal(example.test_variant_int_conversion(10), 10)
	# assert_equal(example.test_variant_int_conversion(10.0), 10)
	# assert_equal(example.test_variant_float_conversion(10.0), 10.0)
	# assert_equal(example.test_variant_float_conversion(10), 10.0)

	# # Test that ptrcalls from GDExtension to the engine are correctly encoding Object and RefCounted.
	# var new_node = Node.new()
	# example.test_add_child(new_node)
	# assert_equal(new_node.get_parent(), example)

	# var new_tileset = TileSet.new()
	# var new_tilemap = TileMap.new()
	# example.test_set_tileset(new_tilemap, new_tileset)
	# assert_equal(new_tilemap.tile_set, new_tileset)
	# new_tilemap.queue_free()

	# # Test variant call.
	# var test_obj = TestClass.new()
	# assert_equal(example.test_variant_call(test_obj), "hello world")

	# Constants.
	assert_equal(Example.FIRST, 0)
	assert_equal(Example.ANSWER_TO_EVERYTHING, 42)
	assert_equal(Example.CONSTANT_WITHOUT_ENUM, 314)

	# BitFields.
	assert_equal(Example.FLAG_ONE, 1)
	assert_equal(Example.FLAG_TWO, 2)
	assert_equal(example.test_bitfield(0), 0)
	assert_equal(example.test_bitfield(Example.FLAG_ONE | Example.FLAG_TWO), 3)

	# Test variant iterator.
	# assert_equal(example.test_variant_iterator([10, 20, 30]), [15, 25, 35])
	# assert_equal(example.test_variant_iterator(null), "iter_init: not valid")

	# RPCs.
	# assert_equal(example.return_last_rpc_arg(), 0)
	# example.test_rpc(42)
	# assert_equal(example.return_last_rpc_arg(), 42)
	# example.test_send_rpc(100)
	# assert_equal(example.return_last_rpc_arg(), 100)

	# Virtual method.
	var event = InputEventKey.new()
	event.key_label = KEY_H
	event.unicode = 72
	get_viewport().push_input(event)
	assert_equal(custom_signal_emitted, ["_input: H", 72])

	# gd extension class calls
	assert_equal(example.test_get_child_node("Label"), example.get_node("Label"))
	example.test_set_position_and_size(Vector2(320, 240), Vector2(100, 200))
	assert_equal(example.get_position(), Vector2(320, 240))
	assert_equal(example.get_size(), Vector2(100, 200))
	# example.test_cast_to()

	# var body = CharacterBody2D.new()
	# var motion = Vector2(1.0, 2.0)
	# body.move_and_collide(motion, true, 0.5, true)
	# example.test_character_body_2d(body)
	# body.queue_free()

	assert_equal(example.test_parent_is_nil(), null)

func _on_Example_custom_signal(signal_name, value):
	custom_signal_emitted = [signal_name, value]
