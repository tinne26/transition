# required triggers
- camera trigger
- switch a flag
- while in area, show tip.
- rotating comment with hint
- interaction touch trigger with hint (press I to interact)

text hints must be drawn behind the player, as a kind of hint block. we
could say we only have one hint at a time, within a delimited area, store
it in level, and once we go out we remove it. we can also check the current
hint alongside an id... ah no, we can set it each frame, no problem.
triggers can also overlap.
hint must stay behind a flag.

the player may also have a hint on itself. basically a talking icon only,
add manually on draw. only "..." or "!", static icons. all hints are
static icons. actually, since sometimes it may have to be placed in a
specific location, see about a "onPlayer" so I can get its pos directly
like that? hmmm, no, better keep a HintBlock available too.
for the text or tips, the Game will work.

trigger messages to main game:
- lock camera to certain position along with player... until the player
  leaves the area
