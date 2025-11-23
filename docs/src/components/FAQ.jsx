import React, { useState } from 'react';
import { ChevronDown } from 'lucide-react';

const FAQ = () => {
    const [openIndex, setOpenIndex] = useState(null);

    const faqs = [
        {
            q: "Does Nanovision run my tests?",
            a: "No. Nanovision acts purely as a reporter and analyzer. You continue to run your tests with your standard language tools (e.g., 'go test', 'dotnet test', 'jest'), and Nanovision consumes the resulting coverage files to generate its reports."
        },
        {
            q: "Is Nanovision free to use?",
            a: "Yes! Nanovision is completely open source and released under the Apache 2.0 license. You can use it freely in both commercial and personal projects without restriction."
        },
        {
            q: "Can I use Nanovision in my CI/CD pipeline?",
            a: "Absolutely. Nanovision is distributed as a single binary with no external dependencies, making it incredibly easy to drop into GitHub Actions, GitLab CI, Jenkins, or any other CI environment."
        }
    ];

    const toggleFAQ = (index) => {
        setOpenIndex(openIndex === index ? null : index);
    };

    return (
        <section id="faq" className="py-24 border-t border-border bg-secondary/20">
            <div className="max-w-4xl mx-auto px-6">
                <div className="text-center mb-16">
                    <h2 className="text-3xl font-bold text-foreground mb-4">FAQ</h2>
                </div>

                <div className="space-y-4">
                    {faqs.map((faq, idx) => (
                        <div
                            key={idx}
                            className="bg-card border border-border rounded-xl overflow-hidden hover:border-primary/50 transition-colors duration-300"
                        >
                            <button
                                className="w-full flex items-center justify-between p-6 text-left focus:outline-none cursor-pointer"
                                onClick={() => toggleFAQ(idx)}
                            >
                                <h3 className="text-lg font-bold text-foreground">{faq.q}</h3>
                                <div className={`transition-transform duration-300 ${openIndex === idx ? 'rotate-180' : ''} text-primary`}>
                                    <ChevronDown size={20} />
                                </div>
                            </button>
                            <div className={`overflow-hidden transition-all duration-300 ease-in-out ${openIndex === idx ? 'max-h-48 opacity-100' : 'max-h-0 opacity-0'}`}>
                                <p className="p-6 pt-0 text-muted-foreground leading-relaxed text-sm">
                                    {faq.a}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default FAQ;
