package lib

func (this *Mlrval) GetType() MVType {
	return this.mvtype
}

func (this *Mlrval) IsError() bool {
	return this.mvtype == MT_ERROR
}

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

func (this *Mlrval) IsVoid() bool {
	return this.mvtype == MT_VOID
}

func (this *Mlrval) IsErrorOrVoid() bool {
	return this.mvtype == MT_ERROR || this.mvtype == MT_VOID
}

func (this *Mlrval) IsBool() bool {
	return this.mvtype == MT_BOOL
}

func (this *Mlrval) GetBoolValue() (boolValue bool, isBoolean bool) {
	if this.mvtype == MT_BOOL {
		return this.boolval, true
	} else {
		return false, false
	}
}

func (this *Mlrval) IsTrue() bool {
	return this.mvtype == MT_BOOL && this.boolval == true
}
func (this *Mlrval) IsFalse() bool {
	return this.mvtype == MT_BOOL && this.boolval == false
}
