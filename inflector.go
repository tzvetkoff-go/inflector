package inflector

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go2c/unidecode"
)

// PluralizationRule represents a regular expression rule for pluralization
type PluralizationRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// SingularizationRule represents a regular expression rule for singularization
type SingularizationRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// IrregularNoun represents a 2-word rule for irregular nouns
type IrregularNoun struct {
	Singular string
	Plural   string
}

// UncountableNoun represents a single word rule for uncountable nouns
type UncountableNoun struct {
	Word string
}

// Acronym represents a lowercase and uppercase acronym rule
type Acronym struct {
	Lower string
	Upper string
}

// Inflector is the inflector inflector type
type Inflector struct {
	PluralizationRules   []*PluralizationRule
	SingularizationRules []*SingularizationRule
	IrregularNouns       []*IrregularNoun
	UncountableNouns     []*UncountableNoun
	AcronymsLower        map[string]string
	AcronymsUpper        map[string]string
	AcronymsUpperList    []string
	AcronymRegexp        *regexp.Regexp
}

// UnderscoreRegexp1 is the first regular expression applied to `Underscore`
// Will transform FOOBar into FOO_Bar
var UnderscoreRegexp1 = regexp.MustCompile("([A-Z]+)([A-Z][a-z])")

// UnderscoreRegexp2 is the second regular expression applied to `Underscore`
// Will transform FooBar into Foo_Bar
var UnderscoreRegexp2 = regexp.MustCompile("([a-z\\d])([A-Z])")

// UnderscoreRegexp3 is the third regular expression applied to `Underscore`
// Will transform spaces and dashes into underscores
var UnderscoreRegexp3 = regexp.MustCompile("[\\s-]+")

// ParameterizeRegexp1 is the first regular expression applied to `Parameterize`
// Will remove all non-word characters
var ParameterizeRegexp1 = regexp.MustCompile("[^\\w]")

// New returns a new Inflector
func New() *Inflector {
	result := &Inflector{}
	result.AcronymsLower = make(map[string]string)
	result.AcronymsUpper = make(map[string]string)
	return result
}

// AddPluralizationRule adds a pluralization rule
func (i *Inflector) AddPluralizationRule(pattern *regexp.Regexp, replacement string) {
	i.PluralizationRules = append(i.PluralizationRules, &PluralizationRule{pattern, replacement})
}

// AddSingularizationRule adds a singularization rule
func (i *Inflector) AddSingularizationRule(pattern *regexp.Regexp, replacement string) {
	i.SingularizationRules = append(i.SingularizationRules, &SingularizationRule{pattern, replacement})
}

// AddIrregularNoun adds an irregular noun
func (i *Inflector) AddIrregularNoun(singular, plural string) {
	i.IrregularNouns = append(i.IrregularNouns, &IrregularNoun{singular, plural})
}

// AddUncountableNoun adds an uncountable noun
func (i *Inflector) AddUncountableNoun(word string) {
	i.UncountableNouns = append(i.UncountableNouns, &UncountableNoun{word})
}

// AddAcronym adds an acronym
func (i *Inflector) AddAcronym(lower, upper string) {
	if _, found := i.AcronymsUpper[upper]; !found {
		i.AcronymsUpperList = append(i.AcronymsUpperList, upper)
	}

	i.AcronymsLower[lower] = upper
	i.AcronymsUpper[upper] = lower

	re := strings.Join(i.AcronymsUpperList, "|")
	i.AcronymRegexp = regexp.MustCompile("(?:\\b)?(" + re + ")(?:\\b)?")
}

// Camelize transforms a word from CamelCase to under_score
func (i *Inflector) Camelize(word string) string {
	words := strings.Split(strings.Replace(strings.ToLower(word), "-", "_", -1), "_")
	camelized := ""

	for _, word := range words {
		if upper, found := i.AcronymsLower[word]; found {
			camelized += upper
			continue
		}

		camelized += strings.Title(word)
	}

	return camelized
}

// Camelize transforms a word from CamelCase to under_score
func Camelize(word string) string {
	return DefaultInflector.Camelize(word)
}

// Underscore transforms a word from under_score to CamelCase
func (i *Inflector) Underscore(word string) string {
	if i.AcronymRegexp != nil {
		word = i.AcronymRegexp.ReplaceAllStringFunc(word, func(match string) string {
			return "_" + i.AcronymsUpper[match]
		})
	}

	word = UnderscoreRegexp1.ReplaceAllString(word, "${1}_${2}")
	word = UnderscoreRegexp2.ReplaceAllString(word, "${1}_${2}")
	word = UnderscoreRegexp3.ReplaceAllString(word, "${1}_${2}")
	word = strings.Replace(word, "__", "_", -1)
	word = strings.Trim(word, "_")
	return strings.ToLower(word)
}

// Underscore transforms a word from under_score to CamelCase
func Underscore(word string) string {
	return DefaultInflector.Underscore(word)
}

// Parameterize transforms a word into a url-friendly-value
func (i *Inflector) Parameterize(word string) string {
	word = i.Transliterate(word)
	word = ParameterizeRegexp1.ReplaceAllString(word, " ")
	word = i.Underscore(word)
	word = strings.Replace(word, "_", "-", -1)
	return word
}

// Parameterize transforms a word into a url-friendly-value
func Parameterize(word string) string {
	return DefaultInflector.Parameterize(word)
}

// Transliterate transliterates non-latin characters
func (i *Inflector) Transliterate(word string) string {
	return unidecode.Unidecode(word)
}

// Transliterate transliterates non-latin characters
func Transliterate(word string) string {
	return DefaultInflector.Transliterate(word)
}

// Pluralize pluralizes an english noun
func (i *Inflector) Pluralize(word string) string {
	lower := strings.ToLower(word)

	for _, rule := range i.UncountableNouns {
		if strings.HasSuffix(lower, rule.Word) {
			return word
		}
	}

	for _, rule := range i.IrregularNouns {
		if strings.HasSuffix(lower, rule.Singular) {
			return word[0:len(word)-len(rule.Singular)] + rule.Plural
		}
		if strings.HasSuffix(lower, rule.Plural) {
			return word
		}
	}

	for _, rule := range i.PluralizationRules {
		if rule.Pattern.MatchString(word) {
			return rule.Pattern.ReplaceAllString(word, rule.Replacement)
		}
	}

	return word + "s"
}

// Pluralize pluralizes an english noun
func Pluralize(word string) string {
	return DefaultInflector.Pluralize(word)
}

// Singularize singularizes an english noun
func (i *Inflector) Singularize(word string) string {
	lower := strings.ToLower(word)

	for _, rule := range i.UncountableNouns {
		if strings.HasSuffix(lower, rule.Word) {
			return word
		}
	}

	for _, rule := range i.IrregularNouns {
		if strings.HasSuffix(lower, rule.Plural) {
			return word[0:len(word)-len(rule.Plural)] + rule.Singular
		}
		if strings.HasSuffix(lower, rule.Singular) {
			return word
		}
	}

	for _, rule := range i.SingularizationRules {
		if rule.Pattern.MatchString(word) {
			return rule.Pattern.ReplaceAllString(word, rule.Replacement)
		}
	}

	return word
}

// Singularize singularizes an english noun
func Singularize(word string) string {
	return DefaultInflector.Singularize(word)
}

// Ordinalize ordinalizes a number
func (i *Inflector) Ordinalize(number int) string {
	if number%100 >= 11 && number%100 <= 13 {
		return strconv.Itoa(number) + "th"
	}

	switch number % 10 {
	case 1:
		return strconv.Itoa(number) + "st"
	case 2:
		return strconv.Itoa(number) + "nd"
	case 3:
		return strconv.Itoa(number) + "rd"
	default:
		return strconv.Itoa(number) + "th"
	}
}

// Ordinalize ordinalizes a number
func Ordinalize(number int) string {
	return DefaultInflector.Ordinalize(number)
}

// DefaultInflector is the default inflector
var DefaultInflector = New()

// Static initialization
func init() {
	// Default pluralization rules for english
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(quiz)$"), "${1}zes")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)^(oxen)$"), "${1}")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)^(ox)$"), "${1}en")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)^(m|l)ice$"), "${1}ice")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)^(m|l)ouse$"), "${1}ice")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(matr|vert|ind)(?:ix|ex)$"), "${1}ices")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(x|ch|ss|sh)$"), "${1}es")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)([^aeiouy]|qu)y$"), "${1}ies")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(hive)$"), "${1}s")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(?:([^f])fe|([lr])f)$"), "${1}${2}ves")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)sis$"), "ses")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)([ti])a$"), "${1}a")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)([ti])um$"), "${1}a")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(buffal|tomat)o$"), "${1}oes")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(bu)s$"), "${1}ses")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(alias|status)$"), "${1}es")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(octop|vir)i$"), "${1}i")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)(octop|vir)us$"), "${1}i")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)^(ax|test)is$"), "${1}es")
	DefaultInflector.AddPluralizationRule(regexp.MustCompile("(?i)s$"), "s")

	// Default singularization rules for english
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(database)s$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(quiz)zes$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(matr)ices$"), "${1}ix")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(vert|ind)ices$"), "${1}ex")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)^(ox)en"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(alias|status)(es)?$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(octop|vir)(us|i)$"), "${1}us")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)^(a)x[ie]s$"), "${1}xis")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(cris|test)(is|es)$"), "${1}is")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(shoe)s$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(o)es$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(bus)(es)?$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)^(m|l)ice$"), "${1}ouse")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(x|ch|ss|sh)es$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(m)ovies$"), "${1}ovie")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(s)eries$"), "${1}eries")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)([^aeiouy]|qu)ies$"), "${1}y")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)([lr])ves$"), "${1}f")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(tive)s$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(hive)s$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)([^f])ves$"), "${1}fe")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(^analy)(sis|ses)$"), "${1}sis")
	DefaultInflector.AddSingularizationRule(
		regexp.MustCompile("(?i)((a)naly|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)(sis|ses)$"), "${1}sis")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)([ti])a$"), "${1}um")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(n)ews$"), "${1}ews")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)(ss)$"), "${1}")
	DefaultInflector.AddSingularizationRule(regexp.MustCompile("(?i)s$"), "")

	// Default irregular nouns for english
	DefaultInflector.AddIrregularNoun("person", "people")
	DefaultInflector.AddIrregularNoun("man", "men")
	DefaultInflector.AddIrregularNoun("woman", "women")
	DefaultInflector.AddIrregularNoun("child", "children")
	DefaultInflector.AddIrregularNoun("tooth", "teeth")
	DefaultInflector.AddIrregularNoun("move", "moves")
	DefaultInflector.AddIrregularNoun("zombie", "zombies")

	// Default uncountable nouns for english
	DefaultInflector.AddUncountableNoun("equipment")
	DefaultInflector.AddUncountableNoun("information")
	DefaultInflector.AddUncountableNoun("rice")
	DefaultInflector.AddUncountableNoun("money")
	DefaultInflector.AddUncountableNoun("species")
	DefaultInflector.AddUncountableNoun("series")
	DefaultInflector.AddUncountableNoun("fish")
	DefaultInflector.AddUncountableNoun("sheep")
	DefaultInflector.AddUncountableNoun("jeans")
	DefaultInflector.AddUncountableNoun("police")
}
