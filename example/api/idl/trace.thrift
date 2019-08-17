namespace go trace

struct Trace {
    1: required string Id,      // id
    2: required string Pid,     // parent_id
    3: required string Cid,     // child_id
}