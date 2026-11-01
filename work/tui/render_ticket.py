"""Render the Sous-Chef ticket to static HTML.

The page's content must survive a browser that blocks scripts: the gallery and
the default preview are written into the file, and the live simulator only
*assembles* what is already there.

This file is the single source of truth for the ticket's content. It writes:

  gallery   the eight variations, as static cards
  default   the default ticket, as static markup inside the live preview
  groups    the line groups the script assembles, as JSON

The script holds no content of its own — it picks a group per control and
concatenates. With scripts blocked, the static default and the static gallery
are what the reader sees.
"""

import json

W = 78
RULE = '─' * W
RULE_IN = ' ' + RULE[:W - 1]

HARNESS = [
    ('opencode', 'opencode', 'opencode'),
    ('pi', 'pi', 'pi --model sonnet'),
    ('custom', 'custom', 'claude --dangerously-skip-permissions'),
]
CREW = {
    3: 'orchestrator, runner, reviewer',
    2: 'runner, and the orchestrator reviews',
    1: 'one agent, and the chef reviews',
}
SOURCES = {
    'ready': ('ready (4)', [
        ('BL05', 'Wire the review queue'),
        ('BL06', 'Sharpen the evidence contract'),
        ('BL11', 'Wire the review queue'),
        ('BL12', 'Sharpen the evidence contract'),
    ]),
    'recipes': ('recipes (6)', [
        ('RCP02', 'Sign off the menu'),
        ('RCP03', 'Wire the review queue'),
        ('RCP04', 'Add the backlog view'),
        ('RCP05', 'Sharpen the evidence contract'),
        ('RCP09', 'Sign off the menu'),
        ('RCP10', 'Wire the review queue'),
    ]),
    'backlog': ('backlog (4)', [
        ('BL01', 'The board needs a way to show evidence'),
        ('BL02', 'A recipe with no ACs is probably an idea'),
        ('BL03', 'The chef wants a weekly digest'),
        ('BL04', 'Trying the board on a real project'),
    ]),
}
BUDGETS = [(2, 8), (3, 12), (5, 24)]
PHASES = ['idle', 'running', 'waiting', 'done']

DEFAULT = dict(harness='opencode', crew=3, source='ready', bi=1, phase='running',
               usedRound=2, usedInter=6, agree=True)

VARIATIONS = [
    ('Three sous-chefs', 'opencode · running', dict(DEFAULT)),
    ('The same, on pi', 'pi · running', dict(DEFAULT, harness='pi')),
    ('A custom harness', 'the command is shown', dict(DEFAULT, harness='custom')),
    ('Two sous-chefs', 'the orchestrator reviews', dict(DEFAULT, crew=2)),
    ('One sous-chef', 'the chef reviews', dict(DEFAULT, crew=1)),
    ('An agent is blocked', 'the orchestrator has the question',
     dict(DEFAULT, source='recipes', phase='waiting')),
    ('The order, unsigned', 'the crew has not started',
     dict(DEFAULT, harness='pi', source='backlog', bi=2, phase='idle', usedRound=0, usedInter=0)),
    ('The build disagrees', 'the stamp says so',
     dict(DEFAULT, phase='done', usedRound=3, usedInter=11, agree=False)),
    ('The round closed', 'the crew is not still working',
     dict(DEFAULT, phase='done', usedRound=3, usedInter=11, agree=True)),
]


def esc(s):
    return str(s).replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;')


def sp(cls, txt):
    return '<span class="%s">%s</span>' % (cls, esc(txt))


def L(parts):
    return ''.join(sp(c, t) for c, t in parts)


def head_lines():
    return [
        L([('c-acc', ' ⋈ maitre · kb')]),
        L([('c-rule', RULE)]),
        L([('c-dim', '  Home    Recipes    Backlog    Ready    Board  '),
           ('c-acc', '▸ Sous-Chef'), ('c-dim', '    Log    Settings')]),
        L([('c-acc', '  THE ORDER')]),
        L([('c-rule', RULE_IN)]),
    ]


def g_harness(hid):
    hcmd = dict((h[0], h[2]) for h in HARNESS)[hid]
    parts = [('c-dim', '  harness  ')]
    for h in HARNESS:
        parts.append(('c-sel', '▸ ' + h[1] + '  ') if h[0] == hid
                     else ('c-faint', '   ' + h[1] + '  '))
    return [L(parts), L([('c-dim', '            '), ('c-faint', 'command: ' + hcmd)])]


def g_crew(n):
    out = [L([('c-dim', '  crew     '), ('c-sel', '▸ ' + str(n) + ' · ' + CREW[n])])]
    for other in sorted(CREW, reverse=True):
        if other != n:
            out.append(L([('c-dim', '           '), ('c-faint', str(other) + ' · ' + CREW[other])]))
    return out


def g_sources(sid):
    parts = [('c-dim', '  sources  ')]
    for key, (label, _items) in SOURCES.items():
        parts.append(('c-sel' if key == sid else 'c-faint',
                      ('▸ ' if key == sid else '  ') + label + '  '))
    return [L(parts)]


def g_budget(bi):
    rounds, inter = BUDGETS[bi]
    return [L([('c-dim', '  budget   '), ('c-dim', '%d rounds · %d interactions' % (rounds, inter))])]


def g_queue(sid, running):
    _label, items = SOURCES[sid]
    out = [L([('c-rule', RULE_IN)]), L([('c-acc', '  THE QUEUE')])]
    for i, (bid, title) in enumerate(items[:4]):
        nxt = 'next' if (running and i == 0) else ''
        out.append(L([('c-dim', '  ' + bid.ljust(6)),
                      ('c-ink' if i == 0 else 'c-dim', title.ljust(42)),
                      ('c-acc', nxt)]))
    if len(items) > 4:
        out.append(L([('c-faint', '  and %d more' % (len(items) - 4))]))
    out.append(L([('c-rule', RULE_IN)]))
    return out


def used_for(bi):
    rounds, inter = BUDGETS[bi]
    return min(2, rounds), max(1, inter // 2)


def g_stamp_head(bi, agree, used_round, used_inter):
    rounds, inter = BUDGETS[bi]
    versions = ('c-dim', 'plan 7 · build 7 · they agree') if agree \
        else ('c-bad', 'plan 7 · build 8 · they disagree')
    return [L([('c-acc', '  THE STAMP'),
               ('c-dim', ' · round %d of %d · %d of %d interactions · ' % (
                   used_round, rounds, used_inter, inter)),
               versions])]


def g_stamp_rows(crew, phase):
    if phase == 'idle':
        return [L([('c-acc', '  THE STAMP'), ('c-faint', ' · the crew has not started')]),
                L([('c-faint', '  press space to start it')])]
    out = []
    if crew == 3:
        out.append(L([('c-good', 'working'), ('c-dim', '  orchestrator  '), ('c-dim', 'routing RCP04')]))
        out.append(L([('c-good', 'working'), ('c-dim', '  runner        '), ('c-dim', 'BL05 · 2 of 3 criteria proven')]))
        out.append(L([('c-dim', 'idle   '), ('c-dim', '  reviewer      '), ('c-dim', 'waiting on the runner')]))
    elif crew == 2:
        out.append(L([('c-good', 'working'), ('c-dim', '  orchestrator  '), ('c-dim', 'routing RCP04 · reviews the runner')]))
        out.append(L([('c-good', 'working'), ('c-dim', '  runner        '), ('c-dim', 'BL05 · 2 of 3 criteria proven')]))
    else:
        out.append(L([('c-good', 'working'), ('c-dim', '  one agent     '), ('c-dim', 'BL05 · 2 of 3 criteria proven')]))
        out.append(L([('c-warn', 'reviews '), ('c-dim', '  the chef      '), ('c-dim', 'the chef holds the review')]))
    if phase == 'waiting':
        out.append(L([('c-warn', 'asked   '), ('c-dim', '  reviewer      '),
                      ('c-warn', 'the scope is ambiguous — the orchestrator has it')]))
    if phase == 'done':
        return [L([('c-good', 'closed  '), ('c-dim', '  the round     '),
                   ('c-dim', '3 of 3 criteria proven · signed off')])]
    return out


def g_footer(phase):
    verb = 'start the crew' if phase in ('idle', 'done') else 'stop the crew'
    return [L([('c-dim', ' space '), ('c-acc', verb),
               ('c-dim', ' · ←→ choose · enter apply · q quit')])]


def render(S):
    """The whole ticket, as lines."""
    used_round, used_inter = used_for(S['bi'])
    if S['phase'] == 'done':
        used_round, used_inter = BUDGETS[S['bi']][0], max(1, BUDGETS[S['bi']][1] - 1)
    if S['phase'] == 'idle':
        used_round, used_inter = 0, 0
    lines = head_lines()
    lines += g_harness(S['harness'])
    lines += g_crew(S['crew'])
    lines += g_sources(S['source'])
    lines += g_budget(S['bi'])
    lines += g_queue(S['source'], S['phase'] == 'running')
    if S['phase'] == 'idle':
        lines += g_stamp_rows(S['crew'], 'idle')
    else:
        lines += g_stamp_head(S['bi'], S['agree'], used_round, used_inter)
        lines += g_stamp_rows(S['crew'], S['phase'])
    lines.append('')
    lines += g_footer(S['phase'])
    return lines


def groups():
    """The line groups the script assembles, keyed by control value."""
    g = {}
    for hid, _l, _c in HARNESS:
        g['harness:' + hid] = g_harness(hid)
    for n in CREW:
        g['crew:%d' % n] = g_crew(n)
    for sid in SOURCES:
        g['sources:' + sid] = g_sources(sid)
        for running in (True, False):
            g['queue:%s:%s' % (sid, 'run' if running else 'idle')] = g_queue(sid, running)
    for bi in range(len(BUDGETS)):
        g['budget:%d' % bi] = g_budget(bi)
        ur, ui = used_for(bi)
        for agree in (True, False):
            g['stamphead:%d:%s' % (bi, 'yes' if agree else 'no')] = g_stamp_head(bi, agree, ur, ui)
            g['stamphead-done:%d:%s' % (bi, 'yes' if agree else 'no')] = g_stamp_head(
                bi, agree, BUDGETS[bi][0], max(1, BUDGETS[bi][1] - 1))
    for crew in CREW:
        for phase in PHASES:
            g['stamprows:%d:%s' % (crew, phase)] = g_stamp_rows(crew, phase)
    for phase in PHASES:
        g['footer:' + phase] = g_footer(phase)
    g['head'] = head_lines()
    return g


def gallery_html():
    cards = []
    for title, cap, st in VARIATIONS:
        shot = '\n'.join(render(st))
        cards.append(
            '<div class="card"><div class="cap"><b>%s</b><span>%s</span></div>'
            '<div class="shot"><pre>%s</pre></div></div>' % (esc(title), esc(cap), shot))
    return '\n'.join(cards)


PAGE = 'work/tui/sous-chef-mockups.html'


def build():
    """Rewrite the page's mockups from this renderer.

    The page keeps three marked regions: the gallery, the default preview, and
    the line groups the script assembles. Everything else is hand-written.
    """
    p = PAGE
    h = open(p).read()

    def splice(text, start_marker, end_marker, body):
        s = text.index(start_marker) + len(start_marker)
        e = text.index(end_marker, s)
        return text[:s] + body + text[e:]

    h = splice(h, '<!--GALLERY-->' + chr(10), chr(10) + '<!--/GALLERY-->', gallery_html())
    h = splice(h, '<!--DEFAULT-->' + chr(10), chr(10) + '<!--/DEFAULT-->', chr(10).join(render(DEFAULT)))
    s = h.index('var G = /*GROUPS*/')
    e = h.index(';', s)
    h = h[:s] + 'var G = /*GROUPS*/' + json.dumps(groups(), ensure_ascii=False) + h[e:]
    open(p, 'w').write(h)
    print('rebuilt', p)


if __name__ == '__main__':
    import sys
    what = sys.argv[1]
    if what == 'gallery':
        print(gallery_html())
    elif what == 'default':
        print('\n'.join(render(DEFAULT)))
    elif what == 'groups':
        print(json.dumps(groups(), ensure_ascii=False))
    elif what == 'build':
        build()
