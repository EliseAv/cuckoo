"""

CUCKOO. Chimes the current time according to a set period.
Inspired by the "Announce the time" setting in macOS, apparently unavailable in Windows.

MIT License

Copyright (c) 2020 Ekevoo.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

"""
from datetime import datetime
from time import sleep
from typing import Iterator

import pyttsx3

PERIOD_IN_MINUTES = 15

WEEKDAYS = tuple("Monday Tuesday Wednesday Thursday Friday Saturday Sunday".split())
ORDINALS = tuple("0 1st 2nd 3rd 4th 5th 6th 7th 8th 9th 10th 11th 12th 13th 14th 15th 16th 17th "
                 "18th 19th 20th 21st 22nd 23rd 24th 25th 26th 27th 28th 29th 30th 31st".split())


def main() -> None:
    """ Announce the time as it happens. """
    for chime in iterate_chimes():
        text = chime_to_text(chime)
        print(chime, text)
        pyttsx3.speak("Hey. " + text)  # windows eats the first few 200ms or so, dunno why.


def iterate_chimes() -> Iterator[datetime]:
    """ Yield every eligible datetime while process is actively running. """
    last_minute = -1
    while True:
        now = datetime.now()
        this_minute = now.hour * 60 + now.minute
        if last_minute != this_minute:
            last_minute = this_minute
            if this_minute % PERIOD_IN_MINUTES == 0:
                yield now
        sleep(7)  # pool somewhat often in order to still be useful when host is suspended


def chime_to_text(chime: datetime) -> str:
    """ Spell the hour, or, if it's midnight, spell the day. """
    hour, minute = chime.hour, chime.minute
    if minute == 0:
        if hour == 0:
            return f"It's {WEEKDAYS[chime.weekday()]} the {ORDINALS[chime.day]}."
        elif hour == 12:
            return "It's noon."
        minute = "hours" if hour != 1 else "hour"
    return f"It's {hour} {minute}."


if __name__ == "__main__":
    main()
