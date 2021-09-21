with per_month_started as (
  select
    substr(Inbetriebnahmedatum, 0, 8) as Monat,
    sum(Nettonennleistung) as Leistung
  from
    EinheitenSolar
  where
    Inbetriebnahmedatum != ''
  group by
    substr(Inbetriebnahmedatum, 0, 8)
),
per_month_retired as (
  select
    substr(DatumEndgueltigeStilllegung, 0, 8) as Monat,
    sum(Nettonennleistung) as Leistung
  from
    EinheitenSolar
  where
    DatumEndgueltigeStilllegung != ''
  group by
    substr(DatumEndgueltigeStilllegung, 0, 8)
),
months as (
  select
    Monat
  from
    per_month_started
  union
  select
    Monat
  from
    per_month_retired
),
per_month_cumulative as (
  select
    m.Monat,
    sum(s.Leistung - r.Leistung) over (
      order by
        m.Monat
    ) as Leistung,
    s.Leistung as Leistung_plus,
    r.Leistung as Leistung_minus
  from
    months m
    left join per_month_started s on (m.Monat = s.Monat)
    left join per_month_retired r on (m.Monat = r.Monat)
)
select
  *
from
  per_month_cumulative
where
  Monat > '2010-00';
