# hpschd

## Mesostics

The Main Text (MT) is simply any text input.

The Spine String (SS) is all Capital Letters (CL).

- 50% Mesostic: The CL is unique between itself and the previous CL.
- 100% Mesostic: The CL is unique between itself, the previous CL, and the next CL.

A "meso-acrostic", arguably another version of a Mesostic, has neither of these limitations.

The algorithm is:

- SS is stored as an addressable array/slice of its characters CL
- Each line of MT is fed into the function
- If a character matches CL, logic happens
- Some kind of line return function (LN) will need to exist to break lines, maybe a setting chopped by word.

### Examples

strings.Contains is probably important. for example with a 50% mesostic:

0. lines LN are defined
1. LN(l-?) is fed to the func
2. function reads each character, knowing what the value of SS(CL-1) is, searching for the first occurance of SS(CL)
3. when it finds the active SS(CL), strings.ToUpper
4. that line is now done, loop moves on to the next SS(CL)


## Display


## I Ching

There are probably dozens if not hundreds of computer programs that simulate the I Ching.

So this doesn't mean to replicate them but to provide a source of randomness for calculating values of the Mesostic that is in line with the kind of approach Cage might do.

For instance, the property of how many words per line could be selected via chance operations.

The SS itself could be chance derived.


## Complexity

How could you demonstrate complexity and chaos here?


## Resources

There are other precedents:

- Nicki Hoffman (python) ::: http://vyh.pythonanywhere.com/psmeso/
- UPenn team (javascript) ::: http://mesostics.sas.upenn.edu/

