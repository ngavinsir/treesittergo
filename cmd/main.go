package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/ngavinsir/treesittergo"
)

//go:embed ts-sql.wasm
var sqlLanguageWasm []byte

//go:embed sql.highlights.scm
var sqlHighlightsQuery string

var q = `select time_id, product
   , last_value(quantity ignore nulls) over (partition by product order by time_id) quantity
   , last_value(quantity respect nulls) over (partition by product order by time_id) quantity
   from ( select times.time_id, product, quantity 
             from inventory partition by  (product) 
                right outer join times on (times.time_id = inventory.time_id) 
   where times.time_id between to_date('01/04/01', 'dd/mm/yy') 
      and to_date('06/04/01', 'dd/mm/yy')) 
   order by  2,1;
select times.time_id, product, quantity from inventory 
   partition by  (product) 
   right outer join times on (times.time_id = inventory.time_id) 
   where times.time_id between to_date('01/04/01', 'dd/mm/yy') 
      and to_date('06/04/01', 'dd/mm/yy') 
   order by  2,1;
select deptno
   , ename
   , hiredate
   , listagg(ename, ',') within group (order by hiredate) over (partition by deptno) as employees
from emp;
 select metric_id ,bsln_guid ,timegroup ,obs_value as obs_value 
 , cume_dist () over (partition by metric_id, bsln_guid, timegroup order by obs_value ) as cume_dist 
 , count(1) over (partition by metric_id, bsln_guid, timegroup ) as n 
 , row_number () over (partition by metric_id, bsln_guid, timegroup order by obs_value) as rrank 
 , percentile_disc(:b7 ) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as mid_tail_value 
 , max(obs_value) over (partition by metric_id, bsln_guid, timegroup ) as max_val 
 , min(obs_value) over (partition by metric_id, bsln_guid, timegroup ) as min_val 
 , avg(obs_value) over (partition by metric_id, bsln_guid, timegroup ) as avg_val 
 , stddev(obs_value) over (partition by metric_id, bsln_guid, timegroup ) as sdev_val 
 , percentile_cont(0.25) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_25 
 , percentile_cont(0.5)  within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_50 
 , percentile_cont(0.75) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_75 
 , percentile_cont(0.90) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_90 
 , percentile_cont(0.95) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_95 
 , percentile_cont(0.99) within group (order by obs_value asc) over (partition by metric_id, bsln_guid, timegroup) as pctile_99
 from timegrouped_rawdata d;
select trim(both ' ' from '  a  ') from dual where trim(:a) is not null;
with
clus_tab as (
select id,
a.attribute_name aname,
a.conditional_operator op,
nvl(a.attribute_str_value,
round(decode(a.attribute_name, n.col,
a.attribute_num_value * n.scale + n.shift,
a.attribute_num_value),4)) val,
a.attribute_support support,
a.attribute_confidence confidence
from table(dbms_data_mining.get_model_details_km('km_sh_clus_sample')) t,
table(t.rule.antecedent) a,
km_sh_sample_norm n
where a.attribute_name = n.col (+) and a.attribute_confidence > 0.55
),
clust as (
select id,
cast(collect(cattr(aname, op, to_char(val), support, confidence)) as cattrs) cl_attrs
from clus_tab
group by id
),
custclus as (
select t.cust_id, s.cluster_id, s.probability
from (select
cust_id
, cluster_set(km_sh_clus_sample, null, 0.2 using *) pset
from km_sh_sample_apply_prepared
where cust_id = 101362) t,
table(t.pset) s
)
select a.probability prob, a.cluster_id cl_id,
b.attr, b.op, b.val, b.supp, b.conf
from custclus a,
(select t.id, c.*
from clust t,
table(t.cl_attrs) c) b
where a.cluster_id = b.id
order by prob desc, cl_id asc, conf desc, attr asc, val asc;`

func main() {
	ctx := context.Background()
	ts, err := treesittergo.New(ctx)
	if err != nil {
		panic(err)
	}
	p, err := ts.NewParser(ctx)
	if err != nil {
		panic(err)
	}
	// defer p.Close(ctx)

	// p.Delete()

	sqlLang, err := ts.LanguageSQL(ctx)
	if err != nil {
		panic(err)
	}
	v, err := p.GetLanguageVersion(ctx, sqlLang)
	if err != nil {
		panic(err)
	}
	log.Printf("sql lang version: %+v\n", v)
	sqlLangName, err := sqlLang.Name(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("lang name: %+v\n", sqlLangName)

	err = p.SetLanguage(ctx, sqlLang)
	if err != nil {
		panic(err)
	}

	tree, err := p.ParseString(ctx, q)
	if err != nil {
		panic(err)
	}
	root, err := tree.RootNode(ctx)
	if err != nil {
		panic(err)
	}
	rootKind, err := root.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node kind: %+v\n", rootKind)
	rootString, err := root.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node string: %+v\n", rootString)
	rootChildCount, err := root.ChildCount(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node child count: %+v\n", rootChildCount)
	child1, err := root.Child(ctx, 0)
	if err != nil {
		panic(err)
	}
	child1Kind, err := child1.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 kind: %+v\n", child1Kind)
	child1String, err := child1.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 string: %+v\n", child1String)
	child1Count, err := child1.ChildCount(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child count: %+v\n", child1Count)
	child1child1, err := child1.Child(ctx, 0)
	if err != nil {
		panic(err)
	}
	child1child1Kind, err := child1child1.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child 1 kind: %+v\n", child1child1Kind)
	child1child1String, err := child1child1.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 1 child 1 string: %+v\n", child1child1String)
	child2, err := root.Child(ctx, 1)
	if err != nil {
		panic(err)
	}
	child2Kind, err := child2.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 2 kind: %+v\n", child2Kind)
	child2String, err := child2.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("child 2 string: %+v\n", child2String)

	q, err := ts.NewQuery(ctx, sqlHighlightsQuery, sqlLang)
	if err != nil {
		panic(err)
	}
	qc, err := ts.NewQueryCursor(ctx)
	if err != nil {
		panic(err)
	}
	qc.Exec(ctx, q, child1child1)
	lastEnd := uint64(0)
	// Iterate over query results
	for {
		m, ok, err := qc.NextMatch(ctx)
		if err != nil {
			panic(err)
		}
		if !ok {
			break
		}
		if m.Captures == nil {
			continue
		}
		for _, c := range m.Captures {
			nodeStartByte, err := c.Node.StartByte(ctx)
			if err != nil {
				panic(err)
			}
			if nodeStartByte < lastEnd {
				continue
			}
			captureName, err := q.CaptureNameForID(ctx, c.ID)
			if err != nil {
				panic(err)
			}
			nodeEndByte, err := c.Node.EndByte(ctx)
			if err != nil {
				panic(err)
			}
			nodeStr, err := c.Node.String(ctx)
			if err != nil {
				panic(err)
			}
			log.Printf("(%d-%d) %s: %s\n", nodeStartByte, nodeEndByte, captureName, nodeStr)
			lastEnd = nodeEndByte
		}
	}
}
