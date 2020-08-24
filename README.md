# hpschd

The Writing-Through Mesostic Generator

A text file for input will be transmogrified into a piece of Mesostic poetry using a configured "Spine String".

## Mesostics

Of course not all mesostics are "writing through" style as Cage did often, they can just as easily be written as they are.

- 50% Mesostic: The CL is unique between itself and the previous CL.
- 100% Mesostic: The CL is unique between itself, the previous CL, and the next CL.
- A "meso-acrostic", arguably another version of a Mesostic, has neither of these limitations.

## Display

REST API currently in development.

The idea is a SpineString and Text are submitted via JSON (?) to the API via POST.

The API calls the mesostic stuff, gets a result, and is responsible for displaying.

This is where chance operations can come into play, e.g. changing typeface and sizes.

## Auto-Display

Not yet in development:

If no input is active, the running app will reach out to a configured endpoint, scrape a (chance-derived?) amount of text, get a randomized SpineString from a list (probably the same list used for go test), and continuously display different mesostics every indeterminate portion of chance derived windows of time.


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

