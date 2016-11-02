# slack channel reorg proposal

# Motivation

I believe we've got too many channels, that are inconsistently named, and too many restrictions on who's allowed to go where.

This increases cognitive load ("where's the right place to say/ask this?") and decreases discoverability ("I had no idea this room was used!")

I'd like to propose:
  - a naming scheme for the 80% use cases
  - curate (and document how/why to curate) channels to fit this naming scheme
  - allow for the other 20% uses cases to remain

Any private channels that I don't know of aren't being addressed here.  I plan on leaving that to the discretion of the private channel owners.

## existing channels

I am able to see 35 existing non-archived channels.  There may be some private channels I'm unable to see.

  - #2017-budget-planning
  - #ada-interns
  - #biz-contractors
  - #climatecontrol-prep
  - #clusters
  - #conferences
  - #deis
  - #dogfooding
  - #fidelity
  - #funfact
  - #general
  - #github-notification
  - #golang
  - #google-contrib
  - #helm
  - #intel
  - #interviews
  - #joyent
  - #korea-advice
  - #kpi-evidence
  - #kpi-reports-2016
  - #kraken-tng
  - #kubecon
  - #kubecon-cpe
  - #newseattleoffice
  - #pipeline
  - #random
  - #samsung-it
  - #samsung-monkey
  - #sdsa-marketing
  - #sdsra-staff
  - #univa-ag-only
  - #visiting-engineers
  - #weekly-summaries
  - #zonar

## prior art

  - https://get.slack.help/hc/en-us/articles/217626408-Organize-and-name-channels
  - http://randsinrepose.com/archives/shaping-your-slack/

# Proposal

## actions

  - archive unused channels
  - rename channels per guidelines
  - create channels per guidelines
  - document guidelines, channel policy, curation policy at doc/slack.md in a follow-on PR

## guidelines

  - long-term channels will have prefixes to denote they're "official"
  - sds-prefixed channels for general chat (contractors, non-cnct folk welcome)
  - cnct-prefixed channels for cnct-specific chat
  - team-prefixed channels for long-term project / engagement teams
  - ping-prefixed channels for notifications
  - office-prefixed channels for administrivia
  - all else is one-off and permitted
    - channels for presentation prep
    - outside collaboration
    - etc
  - address channel proliferation with regular curation
  
## migration

| original channel | destination | why |
| --- | --- | --- |
| #2017-budget-planning | #2017-budget-planning | effort still ongoing |
| #ada-interns | archive | head to #team-kraken | 
| #biz-contractors | archive | no recent activity |
| #climatecontrol-prep | archive | no recent activity |
| #clusters | #cnct-dev | new name |
| #conferences | archive | chat in #cnct-general, channel per conference if we engage |
| #deis | archive | no recent activity, head to #cnct-dev |
| #dogfooding | archive | no recent activity |
| #fidelity | #prep-fidelity | we're still in prep mode for this client |
| #funfact | archive | no recent activity |
| #general | #sds-general | new name |
| #github-notification | #ping-github | new name |
| #golang | archive | head to #cnct-dev, move notifications to ping- or team- channel |
| #google-contrib | #cnct-opensource-meetings | new name |
| #helm | archive | no recent activity, head to #cnct-dev |
| #intel | archive | no recent activity |
| #interviews | #cnct-interviews | only logistics, no candidate feedback |
| #joyent | #team-joyent | new name |
| #korea-advice | archive | no recent activity |
| #kpi-evidence | #cnct-kpi-evidence | new name |
| #kpi-reports-2016 | #kpi-reports-2016 | effort still ongoing |
| #kraken-tng | #team-kraken | long running team |
| #kubecon | #kubecon | we're going to this conference |
| #kubecon-cpe | #proj-kubecon-cpe | effort still ongoing |
| #newseattleoffice | #office-seattle | new name |
| #pipeline | #ping-jenkins | new name |
| #random | #sds-random | new name |
| #samsung-it | #cnct-samsung-it | new name |
| #samsung-monkey | #team-monkey | new name |
| #sdsa-marketing | archive | no recent activity, invite over to #sds-general? |
| #sdsra-staff | #team-sdsra-staff | team involved with sdsra staff meeting |
| #univa-ag-only | archive | no recent activity |
| #visiting-engineers | #cnct-onboarding | how about for all new people vs. interns/visits? |
| #weekly-summaries | archive | no recent activity |
| #zonar | #team-zonar | new name |
| n / a | #office-sanjose | why should seattle have all the fun |
| n / a | #cnct-announce | new channel, keep this tightly groomed, relax our standards for #general and allow it to watercooler |

## results

19 "known" channels

  - #cnct-announce
  - #cnct-opensource-meetings
  - #cnct-dev
  - #cnct-kpi-evidence
  - #cnct-interviews
  - #cnct-onboarding
  - #office-seattle
  - #office-sanjose
  - #ping-github
  - #ping-jenkins
  - #proj-fidelity
  - #proj-kubecon-cpe
  - #sds-general
  - #sds-random
  - #team-sdsra-staff
  - #team-community
  - #team-joyent
  - #team-kraken
  - #team-monkey
  - #team-zonar

4 "one-off" channels

  - #2017-budget-planning
  - #kpi-reports-2016
  - #kubecon
  - #sdsa-marketing

## curation

  - channels will get curated / archived monthly
    - if it's been inactive for the month, and nobody is vouching for it, archive
    - if it's been around for two months, and remains active, find a prefix for it
  - if we are happy with this policy, let's try automating it after a manual attempt or two


