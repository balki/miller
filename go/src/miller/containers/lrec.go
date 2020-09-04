// ================================================================
// This is a hashless implementation of insertion-ordered key-value pairs for
// Miller's fundamental record data structure.
//
// Design:
//
// * It keeps a doubly-linked list of key-value pairs.
//
// * No hash functions are computed when the map is written to or read from.
//
// * Gets are implemented by sequential scan through the list: given a key,
//   the key-value pairs are scanned through until a match is (or is not) found.
//
// * Performance improvement of 10-15% percent over lhmss was found in the C
//   impleentation (for test data).
//
// Motivation:
//
// * The use case for records in Miller is that *all* fields are read from
//   strings & written to strings (split/join), while only *some* fields are
//   operated on.
//
// * Meanwhile there are few repeated accesses to a given record: the
//   access-to-construct ratio is quite low for Miller data records.  Miller
//   instantiates thousands, millions, billions of records (depending on the
//   input data) but accesses each record only once per mapping operation.
//   (This is in contrast to accumulator hashmaps which are repeatedly accessed
//   during a stats run.)
//
// * The hashed impl computes hashsums for *all* fields whether operated on or not,
//   for the benefit of the *few* fields looked up during the mapping operation.
//
// * The hashless impl only keeps string pointers.  Lookups are done at runtime
//   doing prefix search on the key names. Assuming field names are distinct,
//   this is just a few char-ptr accesses which (in experiments) turn out to
//   offer about a 10-15% performance improvement.
//
// * Added benefit: the field-rename operation (preserving field order) becomes
//   trivial.
//
// Notes:
// * nil key is not supported.
// * nil value is not supported.
// ================================================================

package containers

import (
	"bytes"
	"os"

	"miller/lib"
)

// ----------------------------------------------------------------
type Lrec struct {
	FieldCount int64
	Head       *lrecEntry
	Tail       *lrecEntry
}

type lrecEntry struct {
	Key   *string
	Value *lib.Mlrval
	Prev  *lrecEntry
	Next  *lrecEntry
}

// ----------------------------------------------------------------
func NewLrec() *Lrec {
	return &Lrec{
		0,
		nil,
		nil,
	}
}

// ----------------------------------------------------------------
func (this *Lrec) Print() {
	this.Fprint(os.Stdout)
}
func (this *Lrec) Fprint(file *os.File) {
	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := this.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(*pe.Key)
		buffer.WriteString("=")
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("\n")
	(*file).WriteString(buffer.String())
}

// ----------------------------------------------------------------
func newLrecEntry(key *string, value *lib.Mlrval) *lrecEntry {
	kcopy := *key
	vcopy := *value
	return &lrecEntry{
		&kcopy,
		&vcopy,
		nil,
		nil,
	}
}

// ----------------------------------------------------------------
func (this *Lrec) Has(key *string) bool {
	return this.findEntry(key) != nil
}

// ----------------------------------------------------------------
func (this *Lrec) findEntry(key *string) *lrecEntry {
	for pe := this.Head; pe != nil; pe = pe.Next {
		if *pe.Key == *key {
			return pe
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *Lrec) Put(key *string, value *lib.Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newLrecEntry(key, value)
		if this.Head == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = this.Tail
			pe.Next = nil
			this.Tail.Next = pe
			this.Tail = pe
		}
		this.FieldCount++
	} else {
		copy := *value
		pe.Value = &copy
	}
}

// ----------------------------------------------------------------
func (this *Lrec) Prepend(key *string, value *lib.Mlrval) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = newLrecEntry(key, value)
		if this.Tail == nil {
			this.Head = pe
			this.Tail = pe
		} else {
			pe.Prev = nil
			pe.Next = this.Head
			this.Head.Prev = pe
			this.Head = pe
		}
		this.FieldCount++
	} else {
		copy := *value
		pe.Value = &copy
	}
}

// ----------------------------------------------------------------
func (this *Lrec) Get(key *string) *lib.Mlrval {
	pe := this.findEntry(key)
	if pe == nil {
		return nil
	} else {
		return pe.Value
	}
	return nil
}

func (this *Lrec) Clear() {
	this.FieldCount = 0
	// Assuming everything unreferenced is getting GC'ed by the Go runtime
	this.Head = nil
	this.Tail = nil
}

// ----------------------------------------------------------------
func (this *Lrec) Copy() *Lrec {
	that := NewLrec()
	for pe := this.Head; pe != nil; pe = pe.Next {
		that.Put(pe.Key, pe.Value)
	}
	return that
}

// ----------------------------------------------------------------
// Returns true if it was found and removed
func (this *Lrec) Remove(key *string) bool {
	pe := this.findEntry(key)
	if pe == nil {
		return false
	} else {
		this.unlink(pe)
		return true
	}
}

// ----------------------------------------------------------------
func (this *Lrec) unlink(pe *lrecEntry) {
	if pe == this.Head {
		if pe == this.Tail {
			this.Head = nil
			this.Tail = nil
		} else {
			this.Head = pe.Next
			pe.Next.Prev = nil
		}
	} else {
		pe.Prev.Next = pe.Next
		if pe == this.Tail {
			this.Tail = pe.Prev
		} else {
			pe.Next.Prev = pe.Prev
		}
	}
	this.FieldCount--
}

// ----------------------------------------------------------------
//void lrec_prepend(Lrec* prec, char* key, char* value, char free_flags) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//
//	if (pe != NULL) {
//		if (pe->free_flags & FREE_ENTRY_VALUE) {
//			free(pe->value);
//		}
//		pe->value = value;
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//		if (free_flags & FREE_ENTRY_VALUE)
//			pe->free_flags |= FREE_ENTRY_VALUE;
//	} else {
//		pe = mlr_malloc_or_die(sizeof(lrecEntry));
//		pe->key         = key;
//		pe->value       = value;
//		pe->free_flags  = free_flags;
//		pe->quote_flags = 0;
//
//		if (prec->Head == NULL) {
//			pe->Prev   = NULL;
//			pe->Next   = NULL;
//			prec->Head = pe;
//			prec->Tail = pe;
//		} else {
//			pe->Next   = prec->Head;
//			pe->Prev   = NULL;
//			prec->Head->Prev = pe;
//			prec->Head = pe;
//		}
//		prec->field_count++;
//	}
//}

//lrecEntry* lrec_put_after(Lrec* prec, lrecEntry* pd, char* key, char* value, char free_flags) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//
//	if (pe != NULL) { // Overwrite
//		if (pe->free_flags & FREE_ENTRY_VALUE) {
//			free(pe->value);
//		}
//		pe->value = value;
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//		if (free_flags & FREE_ENTRY_VALUE)
//			pe->free_flags |= FREE_ENTRY_VALUE;
//	} else { // Insert after specified entry
//		pe = mlr_malloc_or_die(sizeof(lrecEntry));
//		pe->key         = key;
//		pe->value       = value;
//		pe->free_flags  = free_flags;
//		pe->quote_flags = 0;
//
//		if (pd->Next == NULL) { // Append at end of list
//			pd->Next = pe;
//			pe->Prev = pd;
//			pe->Next = NULL;
//			prec->Tail = pe;
//
//		} else {
//			lrecEntry* pf = pd->Next;
//			pd->Next = pe;
//			pf->Prev = pe;
//			pe->Prev = pd;
//			pe->Next = pf;
//		}
//
//		prec->field_count++;
//	}
//	return pe;
//}

//char* lrec_get_ext(Lrec* prec, char* key, lrecEntry** ppentry) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe != NULL) {
//		*ppentry = pe;
//		return pe->value;
//	} else {
//		*ppentry = NULL;;
//		return NULL;
//	}
//}

//// ----------------------------------------------------------------
//lrecEntry* lrec_get_pair_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	if (position <= 0 || position > prec->field_count) {
//		return NULL;
//	}
//	int sought_index = position - 1;
//	int found_index = 0;
//	lrecEntry* pe = NULL;
//	for (
//		found_index = 0, pe = prec->Head;
//		pe != NULL;
//		found_index++, pe = pe->Next
//	) {
//		if (found_index == sought_index) {
//			return pe;
//		}
//	}
//	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//		MLR_GLOBALS.bargv0, __FILE__, __LINE__);
//	exit(1);
//}

//char* lrec_get_key_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->key;
//	}
//}

//char* lrec_get_value_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->value;
//	}
//}

//// ----------------------------------------------------------------
//void lrec_remove(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//
//	free(pe);
//}

//// ----------------------------------------------------------------
//void lrec_remove_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//
//	free(pe);
//}

// Before:
//   "x" => "3"
//   "y" => "4"  <-- pold
//   "z" => "5"  <-- pnew
//
// Rename y to z
//
// After:
//   "x" => "3"
//   "z" => "4"
//
//void lrec_rename(Lrec* prec, char* old_key, char* new_key, int new_needs_freeing) {
//
//	lrecEntry* pold = lrec_find_entry(prec, old_key);
//	if (pold != NULL) {
//		lrecEntry* pnew = lrec_find_entry(prec, new_key);
//
//		if (pnew == NULL) { // E.g. rename "x" to "y" when "y" is not present
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->key = new_key;
//				if (!new_needs_freeing)
//					pold->free_flags &= ~FREE_ENTRY_KEY;
//			} else {
//				pold->key = new_key;
//				if (new_needs_freeing)
//					pold->free_flags |=  FREE_ENTRY_KEY;
//			}
//
//		} else { // E.g. rename "x" to "y" when "y" is already present
//			if (pnew->free_flags & FREE_ENTRY_VALUE) {
//				free(pnew->value);
//			}
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			}
//			pold->key = new_key;
//			if (new_needs_freeing)
//				pold->free_flags |=  FREE_ENTRY_KEY;
//			else
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			lrec_unlink(prec, pnew);
//			free(pnew);
//		}
//	}
//}

// Cases:
// 1. Rename field at position 3 from "x" to "y when "y" does not exist elsewhere in the srec
// 2. Rename field at position 3 from "x" to "y when "y" does     exist elsewhere in the srec
// Note: position is 1-up not 0-up
//void  lrec_rename_at_position(Lrec* prec, int position, char* new_key, int new_needs_freeing){
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		if (new_needs_freeing) {
//			free(new_key);
//		}
//		return;
//	}
//
//	lrecEntry* pother = lrec_find_entry(prec, new_key);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	pe->key = new_key;
//	if (new_needs_freeing) {
//		pe->free_flags |= FREE_ENTRY_KEY;
//	} else {
//		pe->free_flags &= ~FREE_ENTRY_KEY;
//	}
//	if (pother != NULL) {
//		lrec_unlink(prec, pother);
//		free(pother);
//	}
//}

//// ----------------------------------------------------------------
//void lrec_move_to_head(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_head(prec, pe);
//}

//void lrec_move_to_tail(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_tail(prec, pe);
//}

// ----------------------------------------------------------------
// Simply rename the first (at most) n positions where n is the length of pnames.
//
// Possible complications:
//
// * pnames itself contains duplicates -- we require this as invariant-check
//   from the caller since (for performance) we don't want to check this on every
//   record processed.
//
// * pnames has length less than the current record and one of the new names
//   becomes a clash with an existing name.
//
//   Example:
//   - Input record has names "a,b,c,d,e".
//   - pnames is "d,x,f"
//   - We then construct the invalid "d,x,f,d,e" -- we need to detect and unset
//     the second 'd' field.

//void  lrec_label(Lrec* prec, slls_t* pnames_as_list, hss_t* pnames_as_set) {
//	lrecEntry* pe = prec->Head;
//	sllse_t* pn = pnames_as_list->Head;
//
//	// Process the labels list
//	for ( ; pe != NULL && pn != NULL; pe = pe->Next, pn = pn->Next) {
//		char* new_name = pn->value;
//
//		if (pe->free_flags & FREE_ENTRY_KEY) {
//			free(pe->key);
//		}
//		pe->key = mlr_strdup_or_die(new_name);;
//		pe->free_flags |= FREE_ENTRY_KEY;
//	}
//
//	// Process the remaining fields in the record beyond those affected by the new-labels list
//	for ( ; pe != NULL; ) {
//		char* name = pe->key;
//		if (hss_has(pnames_as_set, name)) {
//			lrecEntry* Next = pe->Next;
//			if (pe->free_flags & FREE_ENTRY_KEY) {
//				free(pe->key);
//			}
//			if (pe->free_flags & FREE_ENTRY_VALUE) {
//				free(pe->value);
//			}
//			lrec_unlink(prec, pe);
//			free(pe);
//			pe = Next;
//		} else {
//			pe = pe->Next;
//		}
//	}
//}

//// ----------------------------------------------------------------
//void lrece_update_value(lrecEntry* pe, char* new_value, int new_needs_freeing) {
//	if (pe == NULL) {
//		return;
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//	pe->value = new_value;
//	if (new_needs_freeing)
//		pe->free_flags |= FREE_ENTRY_VALUE;
//	else
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//}

//// ----------------------------------------------------------------
//static void lrec_link_at_head(Lrec* prec, lrecEntry* pe) {
//
//	if (prec->Head == NULL) {
//		pe->Prev   = NULL;
//		pe->Next   = NULL;
//		prec->Head = pe;
//		prec->Tail = pe;
//	} else {
//		// [b,c,d] + a
//		pe->Prev   = NULL;
//		pe->Next   = prec->Head;
//		prec->Head->Prev = pe;
//		prec->Head = pe;
//	}
//	prec->field_count++;
//}

//static void lrec_link_at_tail(Lrec* prec, lrecEntry* pe) {
//
//	if (prec->Head == NULL) {
//		pe->Prev   = NULL;
//		pe->Next   = NULL;
//		prec->Head = pe;
//		prec->Tail = pe;
//	} else {
//		pe->Prev   = prec->Tail;
//		pe->Next   = NULL;
//		prec->Tail->Next = pe;
//		prec->Tail = pe;
//	}
//	prec->field_count++;
//}

//// ----------------------------------------------------------------
//void lrec_dump(Lrec* prec) {
//	lrec_dump_fp(prec, stdout);
//}

//void lrec_dump_fp(Lrec* prec, FILE* fp) {
//	if (prec == NULL) {
//		fprintf(fp, "NULL\n");
//		return;
//	}
//	fprintf(fp, "field_count = %d\n", prec->field_count);
//	fprintf(fp, "| Head: %16p | Tail %16p\n", prec->Head, prec->Tail);
//	for (lrecEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			pe->key;
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->value == NULL ? "null" :
//			pe->value;
//		fprintf(fp,
//		"| prev: %16p curr: %16p next: %16p | key: %12s | value: %12s |\n",
//			pe->Prev, pe, pe->Next,
//			key_string, value_string);
//	}
//}

//void lrec_dump_titled(char* msg, Lrec* prec) {
//	printf("%s:\n", msg);
//	lrec_dump(prec);
//	printf("\n");
//}

//// ----------------------------------------------------------------
//Lrec* lrec_literal_1(char* k1, char* v1) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_2(char* k1, char* v1, char* k2, char* v2) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	lrec_put(prec, k4, v4, NO_FREE);
//	return prec;
//}

//void lrec_print(Lrec* prec) {
//	FILE* output_stream = stdout;
//	char ors = '\n';
//	char ofs = ',';
//	char ops = '=';
//	if (prec == NULL) {
//		fputs("NULL", output_stream);
//		fputc(ors, output_stream);
//		return;
//	}
//	int nf = 0;
//	for (lrecEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//		if (nf > 0)
//			fputc(ofs, output_stream);
//		fputs(pe->key, output_stream);
//		fputc(ops, output_stream);
//		fputs(pe->value, output_stream);
//		nf++;
//	}
//	fputc(ors, output_stream);
//}

//char* lrec_sprint(Lrec* prec, char* ors, char* ofs, char* ops) {
//	string_builder_t* psb = sb_alloc(SB_ALLOC_LENGTH);
//	if (prec == NULL) {
//		sb_append_string(psb, "NULL");
//	} else {
//		int nf = 0;
//		for (lrecEntry* pe = prec->Head; pe != NULL; pe = pe->Next) {
//			if (nf > 0)
//				sb_append_string(psb, ofs);
//			sb_append_string(psb, pe->key);
//			sb_append_string(psb, ops);
//			sb_append_string(psb, pe->value);
//			nf++;
//		}
//		sb_append_string(psb, ors);
//	}
//	char* rv = sb_finish(psb);
//	sb_free(psb);
//	return rv;
//}
