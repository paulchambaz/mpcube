package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var multiSpaceRe = regexp.MustCompile(`\s{2,}`)

func kid3Available() bool {
	_, err := exec.LookPath("kid3-cli")
	return err == nil
}

func kid3ReadTags(path string) (map[string]string, error) {
	cmd := exec.Command("kid3-cli", "-c", "get all 2", path)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("kid3-cli failed for %s: %w", path, err)
	}
	return parseKid3Output(string(out)), nil
}

func kid3WriteTags(path string, tags map[string]string) error {
	var args []string
	for key, val := range tags {
		args = append(args, "-c", "set "+key+" \""+strings.ReplaceAll(val, "\"", "\\\"")+"\"")
	}
	args = append(args, path)
	cmd := exec.Command("kid3-cli", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("kid3-cli failed for %s: %w: %s", path, err, out)
	}
	return nil
}

func kid3WriteTagsBatch(dir string, tags map[string]string) error {
	var args []string
	for key, val := range tags {
		args = append(args, "-c", "set "+key+" \""+strings.ReplaceAll(val, "\"", "\\\"")+"\"")
	}
	args = append(args, dir)
	cmd := exec.Command("kid3-cli", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("kid3-cli failed for %s: %w: %s", dir, err, out)
	}
	return nil
}

func parseKid3Output(output string) map[string]string {
	tags := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "File:") || strings.HasPrefix(line, "Tag ") {
			continue
		}

		parts := multiSpaceRe.Split(line, 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key != "" && val != "" {
			tags[strings.ToLower(key)] = val
		}
	}
	return tags
}
