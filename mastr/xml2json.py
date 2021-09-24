import json
import xml.parsers.expat

# start-root
# -> start-item-or-end-root
#    -> start-field-or-end-item
#       -> field-value -> end-field -> start-field-or-end-item
#    -> finished
EXPECTING_START_ROOT = 0
EXPECTING_START_ITEM_OR_END_ROOT = 1
EXPECTING_START_FIELD_OR_END_ITEM = 2
EXPECTING_FIELD_VALUE_OR_END_FIELD = 3
FINISHED = 4

class XmlToJson(object):
    def __init__(self, root_element, item_element, callback):
        self._root_element = root_element
        self._item_element = item_element
        self._callback = callback

        self._state = EXPECTING_START_ROOT
        self._item = {}
        self._field_name = None
        self._field_data = ""

    def start_element(self, name, attrs):
        if self._state == EXPECTING_START_ROOT:
            assert name == self._root_element, "expected start of root element"
            self._state = EXPECTING_START_ITEM_OR_END_ROOT
        elif self._state == EXPECTING_START_ITEM_OR_END_ROOT:
            assert name == self._item_element, "expected start of item element"
            self._state = EXPECTING_START_FIELD_OR_END_ITEM 
        elif self._state == EXPECTING_START_FIELD_OR_END_ITEM:
            self._field_name = name
            self._state = EXPECTING_FIELD_VALUE_OR_END_FIELD
        elif self._state == EXPECTING_FIELD_VALUE_OR_END_FIELD:
            raise Error("expected field value or end of field")
        elif self._state == FINISHED:
            raise Error("finished")

    def end_element(self, name):
        if self._state == EXPECTING_START_ROOT:
            raise Error("expected start of root element")
        elif self._state == EXPECTING_START_ITEM_OR_END_ROOT:
            assert name == self._root_element, "expected end of root element"
            self._state = FINISHED
        elif self._state == EXPECTING_START_FIELD_OR_END_ITEM:
            assert name == self._item_element, "expected end of item element"
            self._callback(self._item)
            self._item = {}
            self._state = EXPECTING_START_ITEM_OR_END_ROOT
        elif self._state == EXPECTING_FIELD_VALUE_OR_END_FIELD:
            assert name == self._field_name, "expected end of field"
            self._item[self._field_name] = self._field_data
            self._field_name = None
            self._field_data = ""
            self._state = EXPECTING_START_FIELD_OR_END_ITEM
        elif self._state == FINISHED:
            raise Error("finished")

    def char_data(self, data):
        if self._state != EXPECTING_FIELD_VALUE_OR_END_FIELD:
            return
        self._field_data += data

def main(fname, root_element, item_element):
    p = xml.parsers.expat.ParserCreate(encoding="UTF-16")
    c = XmlToJson(root_element, item_element, lambda x: print(json.dumps(x)))
    p.StartElementHandler = c.start_element
    p.EndElementHandler = c.end_element
    p.CharacterDataHandler = c.char_data

    with open(fname, "rb") as f:
        p.ParseFile(f)

if __name__ == "__main__":
    import sys
    main(sys.argv[1], sys.argv[2], sys.argv[3])
