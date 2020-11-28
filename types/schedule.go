package types

// Schedule Definition of the cost schedule and other parameterizations for wasm vm.
type Schedule struct {
	// Version of the schedule.
	Version uint32 `scale:"version"`

	// PutCodePerByteCost Cost of putting a byte of code into storage.
	PutCodePerByteCost Gas `scale:"put_code_per_byte_cost"`

	// GrowMemCost Gas cost of a growing memory by single page.
	GrowMemCost Gas `scale:"grow_mem_cost"`

	// RegularOpCost Gas cost of a regular operation.
	RegularOpCost Gas `scale:"regular_op_cost"`

	// ReturnDataPerByteCost Gas cost per one byte returned.
	ReturnDataPerByteCost Gas `scale:"return_data_per_byte_cost"`

	// EventDataPerByteCost Gas cost to deposit an event; the per-byte portion.
	EventDataPerByteCost Gas `scale:"event_data_per_byte_cost"`

	// EventPerTopicCost Gas cost to deposit an event; the cost per topic.
	EventPerTopicCost Gas `scale:"event_per_topic_cost"`

	// EventBaseCost Gas cost to deposit an event; the base.
	EventBaseCost Gas `scale:"event_base_cost"`

	// CallBaseCost Base gas cost to call into a contract.
	CallBaseCost Gas `scale:"call_base_cost"`

	// InstantiateBaseCost Base gas cost to instantiate a contract.
	InstantiateBaseCost Gas `scale:"instantiate_base_cost"`

	// DispatchBaseCost Base gas cost to dispatch a runtime call.
	DispatchBaseCost Gas `scale:"dispatch_base_cost"`

	// SandboxDataReadCost Gas cost per one byte read from the sandbox memory.
	SandboxDataReadCost Gas `scale:"sandbox_data_read_cost"`

	// SandboxDataWriteCost Gas cost per one byte written to the sandbox memory.
	SandboxDataWriteCost Gas `scale:"sandbox_data_write_cost"`

	// TransferCost Cost for a simple balance transfer.
	TransferCost Gas `scale:"transfer_cost"`

	// InstantiateCost Cost for instantiating a new contract.
	InstantiateCost Gas `scale:"instantiate_cost"`

	// MaxEventTopics The maximum number of topics supported by an event.
	MaxEventTopics uint32 `scale:"max_event_topics"`

	// MaxStackHeight Maximum allowed stack height.
	//
	// See https://wiki.parity.io/WebAssembly-StackHeight to find out
	// how the stack frame cost is calculated.
	MaxStackHeight uint32 `scale:"max_stack_height"`

	// MaxMemoryPages Maximum number of memory pages allowed for a contract.
	MaxMemoryPages uint32 `scale:"max_memory_pages"`

	// MaxTableSize Maximum allowed size of a declared table.
	MaxTableSize uint32 `scale:"max_table_size"`

	// EnablePrintln Whether the `seal_println` function is allowed to be used contracts.
	// MUST only be enabled for `dev` chains, NOT for production chains
	EnablePrintln bool `scale:"enable_println"`

	// MaxSubjectLen The maximum length of a subject used for PRNG generation.
	MaxSubjectLen uint32 `scale:"max_subject_len"`

	// MaxCodeSize The maximum length of a contract code in bytes. This limit applies to the uninstrumented
	// and pristine form of the code as supplied to `put_code`.
	MaxCodeSize uint32 `scale:"max_code_size"`
}
