package diff

// MaskedChange 构造脱敏字段变更（如凭证值），old/new 使用同一占位符。
func MaskedChange(fieldPath, fieldDesc, mask string, level DiffLevel) DiffRecord {
	return DiffRecord{
		FieldPath: fieldPath,
		FieldDesc: fieldDesc,
		OldValue:  mask,
		NewValue:  mask,
		Level:     level,
	}
}

// AppendRecords 追加变更记录。
func AppendRecords(records []DiffRecord, extra ...DiffRecord) []DiffRecord {
	return append(records, extra...)
}
