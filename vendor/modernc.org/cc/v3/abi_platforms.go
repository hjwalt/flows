package cc

import "encoding/binary"

// abiByteOrders contains byte order information for known architectures.
var (
	abiByteOrders = map[string]binary.ByteOrder{
		"386":     binary.LittleEndian,
		"amd64":   binary.LittleEndian,
		"arm":     binary.LittleEndian,
		"arm64":   binary.LittleEndian,
		"loong64": binary.LittleEndian,
		"ppc64le": binary.LittleEndian,
		"riscv64": binary.LittleEndian,
		"s390x":   binary.BigEndian,
	}

	abiSignedChar = map[[2]string]bool{
		{"freebsd", "arm"}:   false,
		{"freebsd", "arm64"}: false,
		{"linux", "arm"}:     false,
		{"linux", "arm64"}:   false,
		{"linux", "ppc64le"}: false,
		{"linux", "riscv64"}: false,
		{"linux", "s390x"}:   false,
		{"netbsd", "arm"}:    false,
		{"openbsd", "arm64"}: false,

		{"darwin", "amd64"}:  true,
		{"darwin", "arm64"}:  true,
		{"freebsd", "386"}:   true,
		{"freebsd", "amd64"}: true,
		{"illumos", "amd64"}: true,
		{"linux", "386"}:     true,
		{"linux", "amd64"}:   true,
		{"linux", "loong64"}: true,
		{"netbsd", "386"}:    true,
		{"netbsd", "amd64"}:  true,
		{"openbsd", "386"}:   true,
		{"openbsd", "amd64"}: true,
		{"windows", "386"}:   true,
		{"windows", "amd64"}: true,
		{"windows", "arm64"}: true,
	}
)

// abiTypes contains size and alignment information for known OS/arch pairs.
//
// The content is generated by ./cmd/cabi/main.c.
var abiTypes = map[[2]string]map[Kind]ABIType{
	// Linux, generated by GCC 8.3.0
	{"linux", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	{"linux", "386"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 4, 4},
		ULongLong:  {8, 4, 4},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 4, 4},
		LongDouble: {12, 4, 4},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 4, 4},
		UInt64:     {8, 4, 4},
		Float32:    {4, 4, 4},
		Float32x:   {8, 4, 4},
		Float64:    {8, 4, 4},
		Float64x:   {12, 4, 4},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	{"linux", "arm"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {8, 8, 8},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
	},
	{"linux", "arm64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// $ x86_64-w64-mingw32-gcc main.c && wine a.exe
	{"windows", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	// clang version 14.0.0 (https://github.com/llvm/llvm-project.git 329fda39c507e8740978d10458451dcdb21563be)
	// Target: aarch64-w64-windows-gnu
	{"windows", "arm64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {8, 8, 8},
	},
	// $ i686-w64-mingw32-gcc main.c && wine a.exe
	{"windows", "386"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {12, 4, 4},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {12, 4, 4},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	{"darwin", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	{"darwin", "arm64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {8, 8, 8},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// gcc (SUSE Linux) 7.5.0
	{"linux", "s390x"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 8, 8},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 8, 8},
		UInt128:    {16, 8, 8},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 8, 8},
		Float128:   {16, 8, 8},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 8, 8},
	},
	// gcc (FreeBSD Ports Collection) 10.3.0
	{"freebsd", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// gcc (FreeBSD Ports Collection) 11.3.0
	{"freebsd", "arm64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// gcc (FreeBSD Ports Collection) 10.3.0
	{"freebsd", "386"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 4, 4},
		ULongLong:  {8, 4, 4},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 4, 4},
		LongDouble: {12, 4, 4},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 4, 4},
		UInt64:     {8, 4, 4},
		Float32:    {4, 4, 4},
		Float32x:   {8, 4, 4},
		Float64:    {8, 4, 4},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
	},
	// gcc (FreeBSD Ports Collection) 11.3.0
	{"freebsd", "arm"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {8, 8, 8},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
	},
	// gcc (GCC) 8.4.0
	{"openbsd", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
	},
	// OpenBSD clang version 13.0.0
	{"openbsd", "arm64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// OpenBSD clang version 13.0.0
	{"openbsd", "386"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 4, 4},
		ULongLong:  {8, 4, 4},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 4, 4},
		LongDouble: {12, 4, 4},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 4, 4},
		UInt64:     {8, 4, 4},
	},
	// gcc (GCC) 10.3.0
	{"netbsd", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
	},
	// gcc (nb4 20200810) 7.5.0
	{"netbsd", "arm"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {8, 8, 8},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
	},
	// gcc (nb4 20200810) 7.5.0
	{"netbsd", "386"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {4, 4, 4},
		ULong:      {4, 4, 4},
		LongLong:   {8, 4, 4},
		ULongLong:  {8, 4, 4},
		Ptr:        {4, 4, 4},
		Function:   {4, 4, 4},
		Float:      {4, 4, 4},
		Double:     {8, 4, 4},
		LongDouble: {12, 4, 4},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 4, 4},
		UInt64:     {8, 4, 4},
		Float32:    {4, 4, 4},
		Float32x:   {8, 4, 4},
		Float64:    {8, 4, 4},
		Float64x:   {12, 4, 4},
		Float128:   {16, 16, 16},
	},
	// gcc (Ubuntu 11.2.0-7ubuntu2) 11.2.0
	{"linux", "riscv64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
	},
	// gcc (Debian 10.2.1-6) 10.2.1 20210110
	{"linux", "ppc64le"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	// gcc (Loongnix 8.3.0-6.lnd.vec.33) 8.3.0
	{"linux", "loong64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
		Decimal32:  {4, 4, 4},
		Decimal64:  {8, 8, 8},
		Decimal128: {16, 16, 16},
	},
	// gcc (OmniOS 151044/12.2.0-il-0) 12.2.0
	{"illumos", "amd64"}: {
		Void:       {1, 1, 1},
		Bool:       {1, 1, 1},
		Char:       {1, 1, 1},
		SChar:      {1, 1, 1},
		UChar:      {1, 1, 1},
		Short:      {2, 2, 2},
		UShort:     {2, 2, 2},
		Enum:       {4, 4, 4},
		Int:        {4, 4, 4},
		UInt:       {4, 4, 4},
		Long:       {8, 8, 8},
		ULong:      {8, 8, 8},
		LongLong:   {8, 8, 8},
		ULongLong:  {8, 8, 8},
		Ptr:        {8, 8, 8},
		Function:   {8, 8, 8},
		Float:      {4, 4, 4},
		Double:     {8, 8, 8},
		LongDouble: {16, 16, 16},
		Int8:       {1, 1, 1},
		UInt8:      {1, 1, 1},
		Int16:      {2, 2, 2},
		UInt16:     {2, 2, 2},
		Int32:      {4, 4, 4},
		UInt32:     {4, 4, 4},
		Int64:      {8, 8, 8},
		UInt64:     {8, 8, 8},
		Int128:     {16, 16, 16},
		UInt128:    {16, 16, 16},
		Float32:    {4, 4, 4},
		Float32x:   {8, 8, 8},
		Float64:    {8, 8, 8},
		Float64x:   {16, 16, 16},
		Float128:   {16, 16, 16},
	},
}
